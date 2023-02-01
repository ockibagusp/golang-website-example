dep:
	go mod tidy

test:
	go test ./...

test-verbose:
	go test -v ./...

cover:
	go tool cover

cover-html:
	go tool cover -html=coverage.out -o cover.html

cover-show:
	go tool cover -html=coverage.out

run:
	go run app/main/main.go

build:
	go build app/main/main.go