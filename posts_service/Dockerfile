FROM golang:alpine AS builder

WORKDIR /posts_service
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY posts_service/ ./posts_service/
COPY .env ./
COPY protos/ ./protos/
RUN go build -o posts-service ./posts_service/cmd/main.go
CMD ["./posts-service"]