# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS builder

WORKDIR /workspace
COPY go.mod ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 go build -o identity ./cmd/server

FROM gcr.io/distroless/static
ENV HTTP_ADDR=:8080
COPY --from=builder /workspace/identity /identity
ENTRYPOINT ["/identity"] 