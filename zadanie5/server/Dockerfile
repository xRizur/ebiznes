FROM golang:1.21-alpine

WORKDIR /app

RUN apk add --no-cache gcc musl-dev && \
    COPY go.mod go.sum ./ && \
    go mod download && \
    go mod tidy

COPY . .

RUN go build -o main . && chmod +x ./main

CMD ["./main"]

EXPOSE 8080
