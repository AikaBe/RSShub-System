FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o rsshub main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/rsshub .
RUN chmod +x /app/rsshub

EXPOSE 8080

ENTRYPOINT ["./rsshub"]
