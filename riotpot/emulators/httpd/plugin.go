package httpd

import (
	"riotpot/settings"
	"riotpot/utils/emulator"
)

var emu emulator.Emulator

func init() {
	Name := "httpd"
}

func Register() {
	emu := emulator.Emulator{
		Name: "httpd",
	}

	settings.EMULATORS.Register(emu)
}
