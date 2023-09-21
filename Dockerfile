FROM golang:1.20-alpine AS builder
WORKDIR /build
ENV CGO_ENABLED=0
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o story_api

FROM alpine:3.18 AS prod
WORKDIR /app
COPY --from=builder /build/story_api .
USER nobody:nobody

ENTRYPOINT ["/app/story_api", "-c", "config.yml"]
