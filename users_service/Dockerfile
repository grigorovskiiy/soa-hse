FROM golang:alpine AS builder

WORKDIR /users_service
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY users_service/ ./users_service/
COPY .env ./
COPY protos/ ./protos/
RUN go build -o users-service ./users_service/cmd/main.go
CMD ["./users-service"]