FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o paxos-web main.go
EXPOSE 8080
CMD ["./paxos-web"]
