FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o user-service .

FROM gcr.io/distroless/base

COPY --from=builder /app/user-service .

CMD ["/user-service"]
