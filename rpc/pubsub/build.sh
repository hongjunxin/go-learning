#!/bin/bash

go build -o client.out client.go pubsub.pb.go
go build -o server.out server.go pubsub.pb.go
