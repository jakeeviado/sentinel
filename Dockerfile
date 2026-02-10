FROM golang:1.21-bullseye AS builder

RUN apt-get update && apt-get install -y \
    git \
    gcc \
    g++ \
    make \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o sentinel .

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    git \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /build/sentinel /usr/local/bin/sentinel

RUN groupadd -r sentinel && useradd -r -g sentinel sentinel
USER sentinel

ENTRYPOINT ["sentinel"]
CMD ["--help"]
