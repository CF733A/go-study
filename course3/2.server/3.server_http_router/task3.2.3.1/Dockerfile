FROM golang:1.22.2-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /server

FROM alpine:3.18

COPY --from=builder /server /server

COPY .env ./

EXPOSE 8080

CMD ["/server"]