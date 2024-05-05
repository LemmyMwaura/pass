build:
	@go build -o .bin/cmd

start-app:
	@go build -o .bin/cmd && .bin/cmd start

run:
	@.bin/cmd start
