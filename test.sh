#!/bin/sh

curl localhost:8080/api/create/a/1
curl localhost:8080/api/create/b/1
curl localhost:8080/api/create/c/1
curl localhost:8080/api/create/d/1
curl localhost:8080/api/create/e/1

curl localhost:8080/api/delete/e
curl localhost:8080/api/delete/a
curl -m 1 localhost:8080/api/sse/
