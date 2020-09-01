package main

import (
	"os"
	"riotpot/internal/app/riotpot/httpd"
)

func main() {

	args := os.Args
	if args[1] == "--all" {
		httpd.HttpServer()
	}

}
