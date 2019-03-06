# Go microservice

A simple golang microservice with minimal json config interface. 

## Usage 

```bash
# cd to project directory and build executable
$ go build -o bin/microservice .

```

## Docker build

```bash
docker build -t <your-registry-id>/ace-golang:1.11.0 .

```

## Curl timing usage
```
curl -w "@curl-timing.txt" -o /dev/null -s "http://site-to-test

```

## Executing tests
```bash
# clear the cache - this is optional
go clean -testcache
GOCACHE=off go test -v config.go config_test.go schema.go handlers.go middleware.go middleware_test.go handlers_test.go -coverprofile tests/results/cover.out
go tool cover -html=tests/results/cover.out -o tests/results/cover.html


```
## Testing container 
```bash

# start the container
# curl the endpoints
curl -k -H 'Token: MYK*$#GTE457yRH45yhbes%' -w '@curl-timing.txt'  http://127.0.0.1:9000/sys/info/isalive

curl -k -H 'Token: MYK*$#GTE457yRH45yhbes%' -w '@curl-timing.txt' -d '{"username":"MTBLUlVOTkVSQFRBTEtUQUxLLk5FVA==","password":"TFMxNyA5QVQ="}' http://127.0.0.1:9000/login

curl -k -H 'Token: MYK*$#GTE457yRH45yhbes%' -w '@curl-timing.txt' -d '{"apitoken":"bhlcgg7170hpu9f52u7g"}' http://127.0.0.1:9000/alldata
```
