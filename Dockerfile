FROM --platform=${BUILDPLATFORM:-linux/amd64} golang as build-env
LABEL org.opencontainers.image.source=https://github.com/thatnerdjosh/go-example-webservices

WORKDIR /app
COPY . .
ARG TARGETOS TARGETARCH
ENV GOOS $TARGETOS
ENV GOARCH $TARGETARCH
RUN CGO_ENABLED=0 \
    go build -ldflags="-s -w" -o task-webhook ./cmd/task-webhook/main.go

# Create image
FROM scratch
COPY --from=build-env /app/task-webhook /
