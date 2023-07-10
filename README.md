# Go Webservice Examples
**NOTE**: this README is a work in progress

This repository provides some examples of web microservices built in Golang.

## Features
* Multi-arch docker image build
* Services
  * Task webhook service
  * Auth stub service

## Build/Publish docker image

## Run
Using `docker-compose.yaml`:

```bash
docker-compose up -d --build --force-recreate
```

Using `go`:

1. Create and export `.env` to prevent port conflict
```bash
cp .env.example .env
sed -i 's/# AUTH_API_HOST=.*/AUTH_API_HOST=127.0.0.1/' .env
sed -i 's/# AUTH_API_PORT=.*/AUTH_API_PORT=8081/' .env
export $(cat .env | xargs)
```

2. Run services
```bash
nohup go run cmd/auth/main.go &> auth.log &
nohup go run cmd/task-webhook/main.go &> task-webhook.log &
```

3. Cleanup
```bash
kill $(jobs -p)
```

## Known Limitations
* When task-webhook is run in a container it can not restart itself. This would require external orchestration to work properly.
  * Potential workaround would be to set `DOCKER_HOST`, but if it results in restarting itself it will remain stopped.
* Currently all the microservices are bundled together into one image, these should be separated out.
* Logging needs improvement/translation layer implemented.
