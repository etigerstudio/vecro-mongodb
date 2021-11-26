FROM alexeiled/stress-ng AS stress-ng

FROM alpine:3.14
COPY ben-base /
COPY --from=stress-ng /stress-ng / # Grab stress-ng binary from alexeiled/stress-ng
WORKDIR /
ENTRYPOINT ["./ben-base"]