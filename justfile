default:
  just --list
  
build: format vet
  go build -v

test: build
  go test -v ./...

format:
  go fmt

vet:
  go vet

mod:
  go mod download
  go mod tidy
