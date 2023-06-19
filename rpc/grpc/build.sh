#!/bin/bash

go build -o client.out client.go hello.pb.go
go build -o server.out server.go hello.pb.go
