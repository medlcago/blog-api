FROM golang:1.24.5-alpine as builder
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o /app/bin/main ./cmd/main.go

FROM alpine:3
WORKDIR /app
COPY --from=builder /app/bin/main /app/bin/main
ENTRYPOINT ["/app/bin/main"]