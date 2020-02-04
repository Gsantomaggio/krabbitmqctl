# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif


build:
	go build -o bin/krabbitmqctl


release:
	rm -rf bin
	env GOOS=windows GOARCH=amd64 go build -o bin/krabbitmqctl.exe
	env GOOS=darwin GOARCH=amd64 go build  -o bin/krabbitmqctl_darwin_amd64
	env GOOS=linux GOARCH=amd64 go build   -o bin/krabbitmqctl_linux_amd64

