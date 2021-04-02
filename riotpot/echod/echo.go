package main

import (
	"riotpot/services"
)

var Name string = "Echod"

func Echod() services.Service {
	return &Echo{}
}

type Echo struct {
	id   int
	name string
}

func (e *Echo) Init(map[string]interface{}) {}

func (e *Echo) Run() error {
	var err error
	return err
}

func (e *Echo) Stop() error {
	var err error
	return err
}

func (e *Echo) Restart() error {
	var err error
	return err
}

func (e *Echo) Status() error {
	var err error
	return err
}

func (e *Echo) Logger(ch chan<- error) (services.Logger, error) {
	var (
		logger services.Logger
		err    error
	)
	return logger, err
}
