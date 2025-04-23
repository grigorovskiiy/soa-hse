FROM golang:alpine AS builder

WORKDIR /api_gateway_service
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY api_gateway_service/ ./api_gateway_service/
COPY .env ./
COPY protos/ ./protos/
RUN go build -o api-gateway-service ./api_gateway_service/cmd/main.go
CMD ["./api-gateway-service"]