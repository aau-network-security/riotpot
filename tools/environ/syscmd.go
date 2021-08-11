/*
Package environ provides functions used to interact with the environment
*/
package environ

import (
	"fmt"
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
func GetPath(service string) (servicePath string, exists bool) {
	servicePath, err := exec.LookPath(service)

	if err != nil {
		fmt.Printf("Error in getting service %q not found %q	\n", service, err )
		exists = false
		return
	}
	exists = true
	return
}

func ExecuteBackgroundCmd(exec1 string, serviceUsed string) string {
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
	cmd := exec.Command("glider", "-config", "demo.config")
	// cmd.Stdout = os.Stdout
	// cmd.Stdout = &b
	err := cmd.Start()
    if err != nil {
        log.Fatalf("cmd.Run() failed with %s\n", err)
    }
    fmt.Printf("Just ran subprocess %d, exiting\n", cmd.Process.Pid)

	return ""
}

// Execute command and return the output, if any
func ExecuteCmd(app string, args ...string) (output string) {
	cmd := exec.Command(app, args...)
    out, err := cmd.CombinedOutput()
    if err != nil {
        log.Fatalf("cmd.Run() failed with %s\n", err)
    }
    // fmt.Println(strings.Fields(string(out)))
    // fmt.Println(arrays.Contains(arrays.StringToArray(string(out)), "docker-test *"))	
    return string(out)
}

// Execute command and return the output, if any
func CheckDockerExists(name string) (bool) {
	arg := "name="+name
	cmd := ExecuteCmd("docker", "ps", "-a", "--filter", arg)
	// Convert the command output to array and check if name exists
    return arrays.Contains(arrays.StringToArray(cmd), name)
}
