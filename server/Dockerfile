FROM golang:1.22.3-alpine3.19

WORKDIR /home/app

EXPOSE 3000

COPY go.* ./

RUN go mod download

COPY . .

RUN go build -v -o ./bin/server

CMD ["./bin/server"]

