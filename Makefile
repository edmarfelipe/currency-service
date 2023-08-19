run:
	go run -race cmd/api.go

install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GOPATH_DIR}/bin v1.54.1
	golangci-lint --version

build:
	go build -race -o bin/currency_service cmd/api.go

lint:
	golangci-lint run ./...

test:
	go test -race -coverprofile=coverage.txt ./...

