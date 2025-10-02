FROM golang:1.25-alpine AS builder
ENV CGO_ENABLED=0 GO111MODULE=on
WORKDIR /src
RUN apk add --no-cache tzdata ca-certificates git build-base
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev COMMIT=none BUILD_DATE
RUN test -n "$BUILD_DATE" || BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ); \
    GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.Version=$VERSION -X main.Commit=$COMMIT -X main.BuildDate=$BUILD_DATE" -o /out/practice-go ./main.go

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata curl
WORKDIR /app
COPY --from=builder /out/practice-go .
ENV GIN_MODE=release LOG_LEVEL=info PORT=8080
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 CMD curl -fsS http://127.0.0.1:8080/healthz || exit 1
USER nobody
ENTRYPOINT ["./practice-go"]
