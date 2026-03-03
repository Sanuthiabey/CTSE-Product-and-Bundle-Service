FROM golang:1.25-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o product-service cmd/server/main.go

EXPOSE 8080

CMD ["./product-service"]