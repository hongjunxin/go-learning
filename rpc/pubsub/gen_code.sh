#!/bin/bash
protoc --go_out=plugins=grpc:. pubsub.proto
