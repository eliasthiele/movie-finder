FROM golang:1.27.2-alpine AS builder
WORKDIR /app
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
COPY static/ ./static/
EXPOSE 9001
CMD ["./server"]
