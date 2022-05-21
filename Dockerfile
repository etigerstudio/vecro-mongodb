FROM alpine:3.14
COPY ben-mongodb /
WORKDIR /
ENTRYPOINT ["./ben-mongodb"]