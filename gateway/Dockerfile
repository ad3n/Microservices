FROM golang:alpine

RUN mkdir -p /etc/krakend/plugins
RUN apk update
RUN apk add gcc libc-dev bash
WORKDIR /go/src/app
COPY ./gateway .
RUN go run config/processor.go
RUN chmod +x ./builder.sh
RUN ./builder.sh
RUN mv authenticator.so /etc/krakend/plugins
RUN mv api-gateway /usr/bin/api-gateway
RUN chmod +x ./start.sh
RUN apk del gcc libc-dev

CMD ["./start.sh"]

EXPOSE 2727
