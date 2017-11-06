#!/usr/bin/env bash
set -e

if [[ $# -lt 1 ]]; then
    echo "Usage: build_and_push_nginx.sh -t IMAGE_VERSION" >&2
    exit 1
fi

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/alerts-svc main.go

while getopts "v:p" OPTION
do
	case $OPTION in
		v | --version)
			new_tag=app-${OPTARG}
			echo "Using tag: app-$new_tag"
			docker build -t alerts:"$new_tag" .
			docker tag alerts:"$new_tag" imauld/alerts:"$new_tag"
			echo "Built image: alerts:$new_tag"
			;;
		p | --push)
			echo "Pushing image to DockerHub: alerts:$new_tag"
			docker push imauld/alerts:$new_tag
			;;
	esac
done
