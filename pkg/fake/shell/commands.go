package shell

import (
	"fmt"
	"io"
	"strings"
)

// Parses any other command
func c_default(com string, conn io.Writer) (err error) {
	_, err = fmt.Fprintf(conn, "command not found: %s\n", com)
	return
}

// Parses the `enable` command.
func c_enable(com string, conn io.Writer) (err error) {
	// split the string into 2 parts, after the `enable` string
	sp := strings.SplitAfter(com, "enable ")
	rsp := fmt.Sprintf("%s\n%s\n%s\n",
		"cd",
		"enable",
		"exit",
	)

	if len(sp) > 1 {
		// check if the arguments are a flag or an actual command
		// this doesn't go further than just complaining
		if strings.HasPrefix(sp[1], "-") {
			rsp = "enable: bad option: %s\n"
		} else {
			rsp = "enable: no such hash table element: %s\n"
		}

		rsp = fmt.Sprintf(rsp, sp[1])
	}

	_, err = fmt.Fprintf(conn, rsp)
	return
}

// Parses the exit command
func c_exit(com string, conn io.ReadWriteCloser) (err error) {
	_, err = fmt.Fprintf(conn, "bye\n")
	conn.Close()
	return
}

// parses an attempt to execute a file
func c_exec(com string, conn io.Writer) (err error) {
	_, err = fmt.Fprintf(conn, "no such file or directory: %s\n", com)
	return
}
