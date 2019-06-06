FROM golang:1.12-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates tzdata

RUN rm -rf /front-service

RUN mkdir /front-service

ADD . /front-service/

WORKDIR /front-service

RUN go get github.com/gorilla/mux && \
  go get github.com/lgylgy/mpgscore/api

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o front-service .

RUN chmod +x front-service

ENV PORT=3054
ENV GATE_SERVICE_ADDR=localhost:3053

EXPOSE 3054

CMD ["./front-service"]