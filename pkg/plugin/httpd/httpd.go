package main

import (
	"fmt"
	"net/http"

	"github.com/riotpot/internal/globals"
	lr "github.com/riotpot/internal/logger"
	"github.com/riotpot/internal/services"
)

var Plugin string

const (
	name    = "HTTP"
	network = globals.TCP
	port    = 80
)

func init() {
	Plugin = "Httpd"
}

func Httpd() services.Service {
	mx := services.NewPluginService(name, port, network)

	return &Http{
		mx,
	}
}

type Http struct {
	// Anonymous fields from the mixin
	services.Service
}

func (h *Http) Run() (err error) {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(h.valid))

	srv := &http.Server{
		Addr:    h.GetAddress(),
		Handler: mux,
	}

	go h.serve(srv)

	return
}

func (h *Http) serve(srv *http.Server) {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		lr.Log.Fatal().Err(err)
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
	}

	response := fmt.Sprintf("%s%s", head, body)

	fmt.Fprint(w, response)
}
