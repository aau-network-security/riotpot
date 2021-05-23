package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"sync"

	"github.com/riotpot/pkg/fake/shell"
	"github.com/riotpot/pkg/services"
	"github.com/riotpot/tools/errors"
	"github.com/traetox/pty"
	"golang.org/x/crypto/ssh"
)

var Name string

func init() {
	Name = "Sshd"
}

// Inspiration from: https://github.com/jpillora/sshd-lite/
func Sshd() services.Service {

	mixin := services.MixinService{
		Name:     Name,
		Port:     22,
		Running:  make(chan bool, 1),
		Protocol: "tcp",
	}

	return &SSH{
		mixin,
	}
}

type SSH struct {
	services.MixinService
}

func (s *SSH) Run() (err error) {

	// Pre load the configuration for the ssh server
	config := &ssh.ServerConfig{
		PasswordCallback: s.auth,
	}

	// Add a private key for the connections
	config.AddHostKey(s.PrivateKey())

	// before running, migrate the model that we want to store
	s.Migrate(&SSHConn{})
	s.Migrate(&SSHAuth{})

	// convert the port number to a string that we can use it in the server
	var port = fmt.Sprintf(":%d", s.Port)

	// start a service in the `echo` port
	listener, err := net.Listen(s.Protocol, port)
	errors.Raise(err)
	defer listener.Close()

	// create the channel for stopping the service
	s.StopCh = make(chan int, 1)

	// build a channel stack to receive connections to the service
	s.serve(listener, config)

	// update the status of the service
	s.Running <- true

	// Close the channel for stopping the service
	fmt.Print("[x] Service stopped...\n")
	close(s.StopCh)

	return
}

// Function to authenticate the user into the app
func (s *SSH) auth(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
	// Store the credentials before anything
	s.Store(&SSHAuth{
		User:     c.User(),
		Password: string(pass),
	})

	// Currently we don't really care about the credentials
	// any user will have a successful login, as long as the user
	// uses some credentials at all.
	if c.User() != "" && string(pass) != "" {
		return nil, nil
	}
	return nil, fmt.Errorf("invalid pair of username and password")
}

func (s *SSH) serve(listener net.Listener, config *ssh.ServerConfig) {
	// open an infinite loop to receive connections
	fmt.Printf("[%s] Started listenning for connections in port %d\n", Name, s.Port)
	for {
		// Accept the client connection
		client, err := listener.Accept()
		if err != nil {
			return
		}
		defer client.Close()

		// upgrade the connections to ssh
		sshConn, chans, reqs, err := ssh.NewServerConn(client, config)
		if err != nil {
			fmt.Printf("Failed to handshake (%s)", err)
			continue
		}

		sshItem := NewSshConn(sshConn)

		wg := sync.WaitGroup{}
		wg.Add(1)
		// Discard all global out-of-band Requests
		go ssh.DiscardRequests(reqs)
		// Handle all the channels open by the connection
		s.handleChannels(sshItem, chans)
		wg.Wait()
		sshConn.Close()
		sshConn.Conn.Close()
	}
}

func (s *SSH) handleChannels(sshItem SSHConn, chans <-chan ssh.NewChannel) {
	for {
		select {
		case <-s.StopCh:
			// stop the pool
			fmt.Printf("[x] Stopping %s service...\n", s.Name)
			// update the status of the service
			s.Running <- false
			return
		case conn := <-chans:
			//TODO: this line crashes the app when the connection is lost!!!
			// Fix it ASAP!
			// ☟ ☟ ☟
			go s.handleChannel(sshItem, conn)
			// ☝ ☝ ☝
		}
	}
}

// Handles an SSH session
func (s *SSH) handleChannel(sshItem SSHConn, channel ssh.NewChannel) {
	// Check if the type of the channel is a session
	if t := channel.ChannelType(); t != "session" {
		channel.Reject(ssh.UnknownChannelType, "unknown channel type")
		return
	}

	// Accept the channel creation request
	conn, requests, err := channel.Accept()
	if err != nil {
		fmt.Printf("Could not accept channel (%s)", err)
		return
	}

	// handle out-of-band requests
	// Explanation: https://www.wti.com/pages/using-a-reverse-ssh-connection-for-out-of-band-access
	go s.oob(sshItem, requests, conn)
}

// Out-of-band requests handler
// inspired in: https://github.com/traetox/sshForShits/
func (s *SSH) oob(sshItem SSHConn, requests <-chan *ssh.Request, conn ssh.Channel) {
	for req := range requests {
		switch req.Type {
		case "shell":
			if len(req.Payload) == 0 {
				req.Reply(true, nil)
			} else {
				req.Reply(false, nil)
			}
			// Generally, we would put the fake shell
			// under here
			err := s.attachShell(sshItem, conn)
			if err != nil {
				return
			}

		case "pty-req":
			// Responding 'ok' here will let the client
			// know we have a pty ready for input
			req.Reply(true, nil)
		case "window-change":
			continue // no response
		case "env":
			continue // no response
		default:
			fmt.Printf("unkown request: %s (reply: %v, data: %x)", req.Type, req.WantReply, req.Payload)
		}
	}
}

func (s *SSH) attachShell(sshItem SSHConn, conn ssh.Channel) (err error) {
	// load a unix-like fake shell
	shell := shell.New()
	shell.User = sshItem.User
	shell.Host = "ubuntu"

	f, err := pty.StartFaker(shell)
	if err != nil {
		fmt.Printf("Failed to start faker: %v", err)
	}

	close := func() {
		conn.Close()
		shell.Wait()
		f.Close()
	}

	var once sync.Once
	go func() {
		io.Copy(conn, f)
		once.Do(close)
	}()

	go func() {
		io.Copy(f, conn)
		once.Do(close)
	}()

	go func() {
		for rsp := range shell.RspChan {
			s.save(sshItem, rsp)
		}
	}()

	return
}

func (s *SSH) save(sshItem SSHConn, msg []byte) {
	c := &sshItem
	c.Msg = string(msg)
	s.Store(c)
}

// This method returns a private key signer
func (s *SSH) PrivateKey() (key ssh.Signer) {
	// Read the key from a file (???)
	pKey, err := ioutil.ReadFile("configs/keys/riopot_rsa")
	errors.Raise(err)

	// Gets the signer from a key
	key, err = ssh.ParsePrivateKey(pKey)
	errors.Raise(err)

	return
}

type SSHConn struct {
	User          string
	SessionID     []byte
	ClientVersion []byte
	ServerVersion []byte
	RemoteAddr    string
	LocalAddr     string
	Msg           string

	// Request only
	RequestType string
	Payload     []byte
}

type SSHAuth struct {
	User     string
	Password string
}

func NewSshConn(conn *ssh.ServerConn) SSHConn {
	return SSHConn{
		User:          conn.User(),
		SessionID:     conn.SessionID(),
		ClientVersion: conn.ClientVersion(),
		ServerVersion: conn.ServerVersion(),
		RemoteAddr:    conn.RemoteAddr().String(),
		LocalAddr:     conn.LocalAddr().String(),
		Msg:           "",
		RequestType:   "",
		Payload:       []byte{},
	}
}
