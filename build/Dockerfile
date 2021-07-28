FROM golang:1.16-alpine

ARG service

ENV GO111MODULE=on

WORKDIR /app

COPY . .

RUN go mod download

RUN cd ./cmd/${service} && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/run

CMD [ "/app/run" ]
