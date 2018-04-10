FROM alpine:latest

RUN apk --no-cache add curl

COPY ./test-gserver /usr/local/bin

WORKDIR /usr/local/bin

CMD ["test-gserver"]
