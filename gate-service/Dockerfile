FROM golang:1.12-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates tzdata

RUN rm -rf /gate-service

RUN mkdir /gate-service

ADD . /gate-service/

WORKDIR /gate-service

RUN go get github.com/gorilla/mux && \
  go get github.com/lgylgy/mpgscore/api

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o gate-service .

RUN chmod +x gate-service

ENV PORT=3053
ENV MPG_SERVICE_ADDR=localhost:3052
ENV DB_SERVICE_ADDR=localhost:3051

EXPOSE 3053

CMD ["./gate-service"]