FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/server/main.go

FROM alpine:3.18
WORKDIR /app

COPY --from=builder /app/main .

RUN apk add --no-cache ca-certificates

EXPOSE 8080

ENTRYPOINT ["./main"]
