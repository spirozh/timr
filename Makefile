.PHONY: run test

run:
	go run spirozh/timr/cmd/server

test:
	hurl --test test/hurl/*.hurl
