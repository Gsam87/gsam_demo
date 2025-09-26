# Dockerfile
FROM golang:1.24-alpine AS build

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest

WORKDIR /root/
COPY --from=build /app/main .

EXPOSE 8080

CMD ["./main"]