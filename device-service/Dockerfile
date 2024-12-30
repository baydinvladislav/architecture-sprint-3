FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o device-service .

FROM gcr.io/distroless/base

COPY --from=builder /app/device-service .

CMD ["/device-service"]
