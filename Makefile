.PHONY: run test

run:
	go run spirozh/timr/cmd/server

test:
	go run spirozh/timr_test

hurltest:
	hurl --variable host=http://localhost:$server_port --test test/hurl
