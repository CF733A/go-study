FROM golang:1.24.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /main /main

EXPOSE 8080

CMD ["/main"]