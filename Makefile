build:
	@go build -o .bin/cmd

test:
	@go test -v ./...

run:
	@.bin/cmd
