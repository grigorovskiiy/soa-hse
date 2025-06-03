FROM golang:alpine AS builder

WORKDIR /statistic_service
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY statistic_service/ ./statistic_service/
COPY .env ./
COPY protos/ ./protos/
RUN go build -o statistic-service ./statistic_service/cmd/main.go
CMD ["./statistic-service"]