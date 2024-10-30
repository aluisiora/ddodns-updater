-include .env

default: build

build:
	env GOOS=linux GOARCH=amd64 go build -o ./bin/ddodns-updater-linux-x64 *.go

run:
	DO_TOKEN=${DO_TOKEN} DNS_RECORD=${DNS_RECORD} go run *.go

clean:
	rm -rf ./bin/*

