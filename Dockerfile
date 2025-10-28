FROM golang:1.25.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server .
    
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/server .

RUN apk add --no-cache ca-certificates

EXPOSE 8080

CMD ["./server"]
