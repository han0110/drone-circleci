FROM alpine

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

ADD ./release/drone-circleci /go/bin/drone-circleci

ENTRYPOINT ["/go/bin/drone-circleci"]
