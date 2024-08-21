FROM golang:1.21.1-bullseye as builder

WORKDIR /app

COPY ./src/go.* ./
RUN go mod download

COPY ./src ./
COPY .env ./

RUN go build -v -o server ./


FROM debian:bullseye-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/server /app/server
COPY --from=builder /app/.env /app/.env

CMD ["/app/server"]