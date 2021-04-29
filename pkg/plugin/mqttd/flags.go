package main

type ConnectionFlags struct {
	Username, Password, WillRetain, WillFlag, CleanSession bool
	WillQoS                                                uint8
}
