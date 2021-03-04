#!/bin/sh

# create a test of the current project
go test -c

# run the project
cd bin
./riotpot