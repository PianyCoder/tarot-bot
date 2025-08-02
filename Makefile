.PHONY:
#.SILENT:

#make build, make run
build:
	go build -o ./.bin/bot cmd/main.go
run: build
	./.bin/bot