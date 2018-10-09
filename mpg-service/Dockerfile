FROM golang:latest

RUN rm -rf /mpg-service

RUN mkdir /mpg-service

ADD . /mpg-service/

WORKDIR /mpg-service

RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o mpg-service .

RUN chmod +x mpg-service

CMD ["./mpg-service"]