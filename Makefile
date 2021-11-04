build:
	mkdir -p functions
	GOOS=linux GOARCH=amd64 GO111MODULE=auto GOBIN=${PWD}/functions go get ./...
	GOOS=linux GOARCH=amd64 GO111MODULE=auto GOBIN=${PWD}/functions go install ./...