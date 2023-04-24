FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN GO111MODULE="on" CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd

FROM alpine:latest
WORKDIR /task-backend/microservice2
COPY --from=builder /app/app .

EXPOSE 8081
CMD ["./app"]