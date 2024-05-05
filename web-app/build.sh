#!/bin/sh

if ! [ -x "$(command -v goreleaser)" ]; then
	if ! [ -x "$(command -v go)" ]; then
		echo 'Error: Go is not installed'
		exit 1
	fi

	echo 'Info: Installing goreleaser'
	go install github.com/goreleaser/goreleaser@latest
	echo 'Info: Installed goreleaser'
fi

goreleaser build --rm-dist --snapshot
