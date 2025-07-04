#!/usr/bin/env hurl --test

GET http://localhost:8080/selma
HTTP 200
[Asserts]
body contains "Hello, Selma!"
