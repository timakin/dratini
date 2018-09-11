BINARY=app

bundle: GO111MODULE=on go mod download

build: GO111MODULE=on go build -o app

test: GO111MODULE=on go test ./...