FROM alpine:3.14
COPY vecro-mongodb /
WORKDIR /
ENTRYPOINT ["./vecro-mongodb"]