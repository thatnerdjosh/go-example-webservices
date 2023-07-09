# Go Webservice Examples
**NOTE**: this README is a work in progress

This repository provides some examples of web microservices built in Golang.

## Features
* Multi-arch docker image build
* Services
  * Task webhook service

## Build/Publish docker image

## Run
Using `docker-compose.yaml`:

```bash
docker-compose up -d --build --force-recreate
```

Using `go`:

```bash
go run cmd/task-webhook/main.go
```

## Known Limitations
* Currently all the microservices are bundled together into one image, these should be separated out.
* Authentication is currently mocked out simply looking for the authentication header, this should make a call to an auth service.
* Logging needs improvement/translation layer implemented.
