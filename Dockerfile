FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o main .

FROM alpine:3.22

COPY --from=builder /app/main /app/main

EXPOSE 8080

CMD ["/app/main"]
