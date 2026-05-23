FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/app ./cmd/app

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/app /app/app

EXPOSE 8000

CMD ["/app/app"]