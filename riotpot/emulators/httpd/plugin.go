package main

import "common/emulator"
import "settings/settings"

var emulator emulator.Emulator

func init(){
	Name = "httpd"
}

func Register(){
	emulator := emulator.Emulator{
		name = "httpd"
	}

	settings.EMULATORS.register(emulator)
}