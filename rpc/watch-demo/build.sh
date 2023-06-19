#!/bin/bash

go build -o c client.go interface.go
go build -o s server.go interface.go

