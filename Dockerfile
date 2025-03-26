FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download
COPY . .
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o dogobot .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/dogobot .
ENTRYPOINT [ "/app/dogobot" ]
