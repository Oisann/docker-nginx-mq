FROM golang:1.10-alpine as builder

RUN apk --update add git openssh && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*
RUN go get -u github.com/streadway/amqp

WORKDIR /go/src/github.com/oisann/
COPY main.go main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM nginx:alpine
RUN apk add --no-cache bash gawk sed grep bc coreutils
WORKDIR /root/
COPY --from=builder /go/src/github.com/oisann/main .
COPY ./scripts/ /scripts
RUN chmod +x /scripts/*.sh
ENTRYPOINT []
CMD ["/bin/bash", "/scripts/entrypoint.sh"]
