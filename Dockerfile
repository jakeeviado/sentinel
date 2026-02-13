FROM golang:1.21-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git make build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build \
    -a \
    -installsuffix cgo \
    -ldflags '-extldflags "-static"' \
    -o sentinel .

FROM alpine:latest

RUN apk --no-cache add ca-certificates git

WORKDIR /app

COPY --from=builder /build/sentinel /usr/local/bin/sentinel

RUN addgroup -S sentinel && adduser -S sentinel -G sentinel
USER sentinel

ENTRYPOINT ["sentinel"]
CMD ["--help"]
