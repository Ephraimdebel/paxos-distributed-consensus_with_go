FROM golang:1.25

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o server web-service.go

EXPOSE 8080
CMD ["./server"]
