FROM golang:alpine

RUN apk update
WORKDIR /go/src/app
COPY ./service2 .
RUN go get
RUN go build -o service-service .
RUN mv service-service /usr/bin/service-service
RUN chmod +x /usr/bin/service-service

CMD ["service-service"]

EXPOSE 5555
