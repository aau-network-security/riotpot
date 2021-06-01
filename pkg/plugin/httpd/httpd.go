package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/riotpot/pkg/models"
	"github.com/riotpot/pkg/services"
)

var Name string

func init() {
	Name = "Httpd"
}

func Httpd() services.Service {
	mixin := services.MixinService{
		Name:     Name,
		Port:     8080,
		Protocol: "tcp",
		Running:  make(chan bool, 1),
	}

	return &Http{
		mixin,
	}
}

type Http struct {
	// Anonymous fields from the mixin
	services.MixinService
}

func (h *Http) Run() (err error) {
	h.Migrate(&models.Connection{})

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(h.valid))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", h.Port),
		Handler: mux,
	}

	go h.serve(srv)
	h.Running <- true

	for {
		// block until we get the stop signal
		<-h.StopCh

		// send an interrupt signal
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// shut down the server
		srv.Shutdown(context.Background())
		return
	}
}

func (h *Http) serve(srv *http.Server) {
	fmt.Printf("[%s] Started listenning for connections in port %d\n", Name, h.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen:%+s\n", err)
	}
}

// This function handles connections made to a valid path
func (h *Http) valid(w http.ResponseWriter, req *http.Request) {
	var (
		head, body string
	)

	head = `
	<html lang="en">
	<head>
		<!-- Page title -->
		<title>SCADA Login</title>

		<!-- Meta tags -->
		<meta charset="UTF-8">
		<meta id ="viewport" name="viewport" content="width=device-width, initial-scale=1">

		<!-- CSS -->
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta3/dist/css/bootstrap.min.css" 
		rel="stylesheet" integrity="sha384-eOJMYsd53ii+scO/bJGFsiCZc+5NDVN2yr8+0RDqr0Ql0h+rP48ckxlpbzKgwra6" 
		crossorigin="anonymous">
	</head>
	<body>
		<h1>Login</h1><br>
	`

	body = `
		<div class="container">
			<form method="POST">
				<div class="mb-3 row">
					<label for="username" class="form-label">Username</label>
					<input id="username" name="username" class="form-control" type="text" placeholder="Username">
				</div>
				<div class="mb-3 row">
					<label for="password" class="form-label">Password</label>
					<input id="password" name="password" class="form-control" type="password" placeholder="Password">
				</div>
				<button type="submit">Log In</button>
			</form>
		</div>
	</body>
	</html>
	`

	if req.Method == http.MethodPost {
		errormessage := `
		<div class="alert alert-danger">
			<p>Incorrect username or password.</p>
		</div>
		`
		body = errormessage + body

		// save the request
		h.save(req)
	}

	response := fmt.Sprintf("%s%s", head, body)

	fmt.Fprint(w, response)
}

// This function handles connections made to an invalid path
/* func (h *Http) invalid(w http.ResponseWriter, req *http.Request) { ...} */

/*
func (h *Http) loadHandler(path string, valid bool) {
	if valid {
		http.HandleFunc(path, h.valid)
	} else {
		http.HandleFunc(path, h.invalid)
	}
}
*/

func (h *Http) save(req *http.Request) {
	connection := models.NewConnection()
	connection.LocalAddress = "localhost"
	connection.RemoteAddress = req.RemoteAddr
	connection.Protocol = "TCP"
	connection.Service = "HTTP"
	connection.Incoming = true
	connection.Payload = req.PostForm.Encode()

	h.Store(connection)
}
