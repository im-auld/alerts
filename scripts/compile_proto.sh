#!/usr/bin/env bash
set -e

protoc -I ~/repos/kube/proto ~/repos/kube/proto/alerts.proto --go_out=plugins=grpc:proto
protoc -I ~/repos/kube/proto ~/repos/kube/proto/messenger.proto --go_out=plugins=grpc:proto