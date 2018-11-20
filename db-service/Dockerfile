FROM golang:latest

RUN rm -rf /db-service

RUN mkdir /db-service

ADD . /dbservice/

WORKDIR /db-service

RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o db-service .

RUN chmod +x db-service

CMD ["./db-service"]