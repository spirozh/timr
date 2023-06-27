#!/bin/sh

curl -i localhost:8080/api/create/a/0
curl -i localhost:8080/api/create/b/0
curl -i localhost:8080/api/create/c/0
curl -i localhost:8080/api/create/d/0
curl -i localhost:8080/api/create/e/0

curl -i localhost:8080/api/delete/e
curl -i localhost:8080/api/delete/a
