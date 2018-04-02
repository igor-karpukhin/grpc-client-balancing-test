FROM alpine:latest
MAINTAINER IgorKarpukhin

RUN apk --no-cache add curl

COPY ./test-gserver /usr/local/bin

CMD ["test-gserver"]
