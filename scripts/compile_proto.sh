#!/usr/bin/env bash
set -e

protoc -I $GOPTAH/src/github.com/kube-message/proto $GOPTAH/src/github.com/kube-message/proto/alerts.proto --go_out=plugins=grpc:proto
