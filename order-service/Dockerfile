FROM golang:1.23.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/app/
RUN go build -o main .

FROM alpine:3.21.0

WORKDIR /app

RUN apk add --no-cache libc6-compat

COPY ./internal/templates /app/internal/templates
COPY ./internal/static /app/internal/static
COPY --from=builder /app/cmd/app/main .

ENTRYPOINT ["./main"]
