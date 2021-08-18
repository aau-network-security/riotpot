/*
Package environ provides functions used to interact with the environment
*/
package environ

import (
	"log"
	"os/exec"
	"github.com/riotpot/tools/arrays"
)

// get the full binary path on host for a given service 
func GetPath(service string) (servicePath string) {
	servicePath, err := exec.LookPath(service)

	if err != nil {
		log.Fatalf("Error in getting service %q with error %q \n", service, err )
	}
	return servicePath
}

// Execute terminal command in async mode i.e. in background
func ExecuteBackgroundCmd(app string, args ...string) {
	cmd := exec.Command(app, args...)
	err := cmd.Start()
    
    if err != nil {
        log.Fatalf("cmd.Run() for command %q %q failed with %s\n", app, args, err)
    }
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

// Outputs if docker container of the given name exists already
func CheckDockerExists(name string) (bool) {
	arg := "name="+name
	cmd := ExecuteCmd("docker", "ps", "-a", "--filter", arg)
	// Convert the command output to array and check if name exists
    return arrays.Contains(arrays.StringToArray(cmd), name)
}
