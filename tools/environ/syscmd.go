/*
Package environ provides functions used to interact with the environment
*/
package environ

import (
	// "fmt"
	"os/exec"
	// "os"
	"log"
	// "strings"
	"github.com/riotpot/tools/arrays"
	// "bytes"
)

/*
	Check if the port on the host machine is busy or not
	this is used for plugins to play on the host
*/
func GetPath(service string) (servicePath string) {
	servicePath, err := exec.LookPath(service)

	if err != nil {
		log.Fatalf("Error in getting service %q with error %q \n", service, err )
	}
	return servicePath
}

func ExecuteBackgroundCmd(app string, args ...string) {
	// executable := &exec.Cmd {
	// 	Path: Exec,
	// 	Args: []string{Exec, "version" },
	// 	Stdout: os.Stdout,
	// 	Stderr: os.Stdout,
	// }
	// var b bytes.Buffer
	// executable.Stdout = &b
	// executable.Stderr = &b
	
	// err := executable.Run()
	// fmt.Println(string(b.Bytes()))
	// if err != nil {
	// 	fmt.Printf("Service %q not found	\n", Service)
	// }
	cmd := exec.Command(app, args...)
	// cmd.Stdout = os.Stdout
	// cmd.Stdout = &b
	err := cmd.Start()
    if err != nil {
        log.Fatalf("cmd.Run() for command %q %q failed with %s\n", app, args, err)
    }
    // fmt.Printf("Just ran subprocess %d, exiting\n", cmd.Process.Pid)

	// return ""
}

func ExecuteBackgroundCmd1(app string, args ...string) {
	// executable := &exec.Cmd {
	// 	Path: Exec,
	// 	Args: []string{Exec, "version" },
	// 	Stdout: os.Stdout,
	// 	Stderr: os.Stdout,
	// }
	// var b bytes.Buffer
	// executable.Stdout = &b
	// executable.Stderr = &b
	
	// err := executable.Run()
	// fmt.Println(string(b.Bytes()))
	// if err != nil {
	// 	fmt.Printf("Service %q not found	\n", Service)
	// }
	cmd := exec.Command(app, args...)
	// cmd.Stdout = os.Stdout
	// cmd.Stdout = &b
	err := cmd.Start()
    if err != nil {
        log.Fatalf("cmd.Run() for command %q %q failed with %s\n", app, args, err)
    }
    // fmt.Printf("Just ran subprocess %d, exiting\n", cmd.Process.Pid)

	// return ""
}

// Execute command and return the output, if any
func ExecuteCmd(app string, args ...string) (output string) {
	cmd := exec.Command(app, args...)
    out, err := cmd.CombinedOutput()

    if err != nil {
        log.Fatalf("cmd.Run() for command %q %q failed with %s\n", app, args, err)
    }

    return string(out)
}

// Outputs if docker contianer of the given name exists already
func CheckDockerExists(name string) (bool) {
	arg := "name="+name
	cmd := ExecuteCmd("docker", "ps", "-a", "--filter", arg)
	// Convert the command output to array and check if name exists
    return arrays.Contains(arrays.StringToArray(cmd), name)
}
