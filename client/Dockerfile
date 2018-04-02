FROM alpine:latest
MAINTAINER IgorKarpukhin

RUN apk --no-cache add curl

COPY ./test-gclient /usr/local/bin

CMD ["test-gclient"]
