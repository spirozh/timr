.PHONY: run test

run:
	go run github.com/spirozh/timr/cmd/server

test:
	cd src ; go test ./...

hurltest:
	hurl --variable host=http://localhost:$server_port --test test/hurl
