FROM alpine:3.14
COPY ben-base /
WORKDIR /
ENTRYPOINT ["./ben-base"]