FROM golang:1.14

WORKDIR /auth-go-mongodb
COPY    go.mod  .
COPY go.sum .
RUN go mod download 
COPY ./server ./server
RUN cd server \
   go build
EXPOSE "8084"


CMD ["/auth-go-mongodb/server/server"]

