SHELL=/bin/bash
SRC=$(shell find . -name "*.go")

bin/palette: $(SRC)
	go build -o bin/palette cmd/palette.go
	if [[ ! -x /usr/bin/palette ]]; then sudo ln -s $(shell pwd)/bin/palette /usr/bin; fi