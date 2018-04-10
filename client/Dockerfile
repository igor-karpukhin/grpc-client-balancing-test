FROM alpine:latest

RUN apk --no-cache add curl

COPY ./test-gclient /usr/local/bin
COPY ./docker-entrypoint.sh /

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["test-gclient"]
