// This package implements multiple errors specific to riotpot
package errors

import (
	"fmt"
	"log"
)

const (
	EmulatorNotInstalled = "Emulator not installed."
)

func Error(name string, errorString string) {
	err := fmt.Errorf("[!] Error: %s\n	>	%s", name, errorString)
	fmt.Println(err)
}

// Check if there is an error and throw a fatal
func Raise(err error) {
	if err != nil {
		log.Fatalf("[!] Error: %v", err)
	}
}
