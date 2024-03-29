# build stage
FROM golang:1.20rc1-alpine3.17 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# RUn stage
FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["/app/main"]