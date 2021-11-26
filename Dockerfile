FROM alexeiled/stress-ng AS stress-ng

FROM alpine:3.14
COPY ben-base /
# Grab stress-ng binary from alexeiled/stress-ng
COPY --from=0 /stress-ng /
WORKDIR /
ENTRYPOINT ["./ben-base"]