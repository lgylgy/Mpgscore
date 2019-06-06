FROM golang:1.12-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates tzdata

RUN rm -rf /db-service

RUN mkdir /db-service

ADD . /db-service/

WORKDIR /db-service

RUN go get github.com/gorilla/mux && \
  go get gopkg.in/mgo.v2 && \
  go get github.com/rs/xid && \
  go get github.com/lgylgy/mpgscore/api

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o db-service .

RUN chmod +x db-service

ENV PORT=3051
ENV TLS=true
ENV MONGODB=mongodb://<user>:<password>@cluster0-shard-00-00-2hji4.gcp.mongodb.net:27017,cluster0-shard-00-01-2hji4.gcp.mongodb.net:27017,cluster0-shard-00-02-2hji4.gcp.mongodb.net:27017

EXPOSE 3051

CMD ["./db-service"]