FROM golang:1.12-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates tzdata

RUN rm -rf /mpg-service

RUN mkdir /mpg-service

ADD . /mpg-service/

WORKDIR /mpg-service

RUN go get github.com/gorilla/mux && \
  go get github.com/lgylgy/mpgscore/api

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o mpg-service .

RUN chmod +x mpg-service

ENV PLAYERDB=3052

EXPOSE 3052

CMD ["./mpg-service"]