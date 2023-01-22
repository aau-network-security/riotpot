package main

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/riotpot/internal/globals"
	"github.com/riotpot/internal/logger"
	"github.com/riotpot/internal/plugins"
	"github.com/riotpot/internal/services"
	"github.com/riotpot/pkg/fake/shell"

	"github.com/traetox/pty"
	"golang.org/x/crypto/ssh"
)

var Plugin string

const (
	name    = "SSH"
	network = globals.TCP
	port    = 22
)

func init() {
	Plugin = "Sshd"
}

// Inspiration from: https://github.com/jpillora/sshd-lite/
func Sshd() services.Service {

	mx := services.NewPluginService(name, port, network)
	pKey := plugins.NewPrivateKey(plugins.DefaultKey)
	pem := pKey.GetPEM()

	return &SSH{
		Service:    mx,
		wg:         sync.WaitGroup{},
		privateKey: pem,
	}
}

type SSH struct {
	services.Service
	wg         sync.WaitGroup
	privateKey []byte
}

func (s *SSH) Run() (err error) {

	// Preload the configuration for the ssh server
	config := &ssh.ServerConfig{
		PasswordCallback: s.auth,
	}

	// Add a private key for the connections
	priv := s.PrivateKey()
	config.AddHostKey(priv)

	listener, err := net.Listen(s.GetNetwork().String(), s.GetAddress())
	if err != nil {
		return
	}
	defer listener.Close()

	// build a channel stack to receive connections to the service
	s.serve(listener, config)
	return
}

// Function to authenticate the user into the app
func (s *SSH) auth(c ssh.ConnMetadata, pass []byte) (perms *ssh.Permissions, err error) {
	// Currently we don't really care about the credentials
	// any user will have a successful login, as long as the user
	// uses some credentials at all.
	if c.User() != "" && string(pass) != "" {
		return
	}

	err = fmt.Errorf("invalid pair of username and password")
	return
}

func (s *SSH) serve(listener net.Listener, config *ssh.ServerConfig) {
	// open an infinite loop to receive connections
	for {
		// Accept the client connection
		client, err := listener.Accept()
		if err != nil {
			logger.Log.Error().Err(err)
			continue
		}

		// upgrade the connections to ssh
		sshConn, chans, reqs, err := ssh.NewServerConn(client, config)
		if err != nil {
			logger.Log.Error().Err(err)
			continue
		}

		sshItem := NewSshConn(sshConn)

		// Discard all global out-of-band Requests
		go ssh.DiscardRequests(reqs)
		// Handle all the channels open by the connection
		go s.handleChannels(sshItem, chans)
	}
}

func (s *SSH) handleChannels(sshItem SSHConn, chans <-chan ssh.NewChannel) {
	for conn := range chans {
		go s.handleChannel(sshItem, conn)
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
		logger.Log.Error().Err(err)
		return
	}

	// handle out-of-band requests
	// Explanation: https://www.wti.com/pages/using-a-reverse-ssh-connection-for-out-of-band-access
	go s.oob(sshItem, requests, conn)
}

// Out-of-band requests handler
func (s *SSH) oob(sshItem SSHConn, requests <-chan *ssh.Request, conn ssh.Channel) {

	for req := range requests {

		switch req.Type {
		case "shell":
			if len(req.Payload) > 0 {
				logger.Log.Error().Msgf("Shell command ignored", req.Payload)
			}

			// Give a shell to the client
			err := s.attachShell(sshItem, conn)
			if err != nil {
				logger.Log.Error().Err(err)
			}

			req.Reply(err == nil, nil)
		case "pty-req":
			// Responding 'ok' here will let the client
			// know we have a pty ready for input
			req.Reply(true, nil)
		case "window-change":
			continue // no response
		case "env":
			continue // no response
		default:
			logger.Log.Info().Msgf("Unkown request: %s (reply: %v, data: %x)", req.Type, req.WantReply, req.Payload)
		}
	}
}

func (s *SSH) attachShell(sshItem SSHConn, conn ssh.Channel) (err error) {
	// load a unix-like fake shell
	shell := shell.New(sshItem.User, "ubuntu")

	f, err := pty.StartFaker(shell)
	if err != nil {
		return
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

	return
}

// This method returns a private key signer
func (s *SSH) PrivateKey() (key ssh.Signer) {
	// Gets the signer from a key
	key, err := ssh.ParsePrivateKey(s.privateKey)
	if err != nil {
		logger.Log.Fatal().Err(err)
	}

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
