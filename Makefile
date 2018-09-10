VERSION=0.9.1
TARGETS_NOVENDOR=$(shell glide novendor)

all: bin/dratini bin/dratini_recover

build-cross: cmd/dratini/dratini.go cmd/dratini_recover/dratini_recover.go dratini/*.go
	GOOS=linux GOARCH=amd64 go build -o bin/linux/amd64/dratini-${VERSION}/dratini cmd/dratini/dratini.go
	GOOS=linux GOARCH=amd64 go build -o bin/linux/amd64/dratini-${VERSION}/dratini_recover  cmd/dratini_recover/dratini_recover.go
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin/amd64/dratini-${VERSION}/dratini cmd/dratini/dratini.go
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin/amd64/dratini-${VERSION}/dratini_recover cmd/dratini_recover/dratini_recover.go

dist: build-cross
	cd bin/linux/amd64 && tar zcvf dratini-linux-amd64-${VERSION}.tar.gz dratini-${VERSION}
	cd bin/darwin/amd64 && tar zcvf dratini-darwin-amd64-${VERSION}.tar.gz dratini-${VERSION}

bin/dratini: cmd/dratini/dratini.go dratini/*.go
	go build -o bin/dratini cmd/dratini/dratini.go

bin/dratini_recover: cmd/dratini_recover/dratini_recover.go dratini/*.go
	go build -o bin/dratini_recover cmd/dratini_recover/dratini_recover.go

bin/dratini_client: samples/client.go
	go build -o bin/dratini_client samples/client.go

bundle:
	glide install

fmt:
	@echo $(TARGETS_NOVENDOR) | xargs go fmt

check:
	go test -v $(TARGETS_NOVENDOR)

clean:
	rm -rf bin/*
