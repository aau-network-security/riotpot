package main

import (
	"riotpot/settings"
	"riotpot/utils/emulator"
)

var emu emulator.Emulator
var Name string

func init() {
	Name = "httpd"
}

func Register() {
	emu := emulator.Emulator{
		Name: "httpd",
	}

	settings.EMULATORS.Register(emu)
}
