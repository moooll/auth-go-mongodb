FROM golang:latest

WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download 
COPY . .
WORKDIR /app/server
RUN go build
EXPOSE "8084"


ENTRYPOINT ["./server"]

