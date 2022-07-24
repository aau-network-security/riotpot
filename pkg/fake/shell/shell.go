// This package implements shell
// interfaces, which dictates the necessary
// methods to be present in other shell
// structures
//
// Thanks to https://github.com/traetox/sshForShits/ for the inspiration!
// The ssh shell could have not been done without his approach.
package shell

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
)

type shellface interface {
	prompt()
	command()
	commands()
}

func New(user string, host string) *shell {
	commands := [][]interface{}{
		{"enable", c_enable},
		{"exit", c_exit},
		{"./", c_exec},
	}

	return &shell{
		comm:     commands,
		mu:       &sync.Mutex{},
		User:     user,
		Host:     host,
		Path:     "",
		RspChan:  make(chan []byte, 10),
		doneChan: make(chan error, 2),
		Running:  false,
	}
}

type shell struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	closer io.Closer

	shellface

	User    string
	Host    string
	Path    string
	Running bool

	comm [][]interface{}

	RspChan  chan []byte
	doneChan chan error

	mu *sync.Mutex
}

func (s *shell) SetIo(conn io.ReadWriteCloser) {
	s.stdin = conn
	s.stdout = conn
	s.stderr = conn
	s.closer = conn
}

// Method necessary to start fake shells on ssh
func (s *shell) SetReadWriteCloser(r io.Reader, w io.Writer, c io.Closer) {
	s.stdin = r
	s.stdout = w
	s.stderr = w
	s.closer = c
}

func (s *shell) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Running {
		return errors.New("already running")
	}
	s.Running = true
	go s.terminal()

	return nil
}

func (s *shell) Wait() error {
	s.mu.Lock()
	if !s.Running {
		s.mu.Unlock()
		return errors.New("not running")
	}
	s.mu.Unlock()
	return <-s.doneChan
}

func (s *shell) terminal() {
	br := bufio.NewReader(s.stdin)

	for {
		// sends the prompt string
		fmt.Fprintf(s.stdout, "%s", s.prompt())
		// reads the response
		lineBytes, err := br.ReadBytes('\n')
		if err != nil {
			break
		}

		// send the response to the channel of responses
		s.RspChan <- lineBytes
		line := string(lineBytes)
		// remove the line endings to compare the strings to regular commands
		line = strings.TrimRight(line, "\r\n")

		s.commands(line)
	}

	s.Running = false
	s.doneChan <- s.closer.Close()
}

func (s *shell) prompt() string {
	return fmt.Sprintf("%s@%s:~%s# ", s.User, s.Host, s.Path)
}

func (s *shell) commands(line string) {
	commands := strings.Split(line, ";")

	for _, command := range commands {
		s.command(command)
	}
}

func (s *shell) command(line string) {
	for _, p := range s.comm {

		// The first value is the `command` as a string, stringify it.
		command := fmt.Sprintf("%v", p[0])

		// The second is a function that takes exactly 1 argument,
		// the command.
		// We do this so the function can send more than one message,
		// close the connection, etc.
		fn := reflect.ValueOf(p[1])
		args := make([]reflect.Value, 2)

		// Parses the value as an argument
		args[0] = reflect.ValueOf(line)
		args[1] = reflect.ValueOf(s.stdout)

		// find the captured command that matches the list of commands
		// able to parse.
		if strings.HasPrefix(line, command) {
			fn.Call(args)
			return
		}
	}

	c_default(line, s.stdout)
}
