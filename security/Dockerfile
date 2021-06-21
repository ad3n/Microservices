FROM golang:alpine

RUN apk update
WORKDIR /go/src/app
COPY ./security .
RUN go get
RUN go build -o security-service .
RUN mv security-service /usr/bin/security-service
RUN chmod +x /usr/bin/security-service

CMD ["security-service"]

EXPOSE 7777
