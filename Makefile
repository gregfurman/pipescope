BINARY_NAME=pipescope

build:
	CGO_ENABLED=0 go build -mod=mod -o ./targets/${BINARY_NAME} main.go

mod:
	go mod download

run: build
	./${BINARY_NAME}

clean:
	go clean -testcache
	rm -rf ./targets

test:
	go test -race -v ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go get .

vet:
	go vet

lint:
	golangci-lint run -v