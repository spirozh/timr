.PHONY: run test

run:
	cd src ; go run github.com/spirozh/timr/cmd/server

test:
	cd src ; go test ./...
