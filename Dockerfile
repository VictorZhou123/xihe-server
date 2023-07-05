FROM golang:1.18.10-bullseye as BUILDER

# build binary
COPY . /go/src/github.com/opensourceways/xihe-server
RUN cd /go/src/github.com/opensourceways/xihe-server && GO111MODULE=on CGO_ENABLED=0 go build

# copy binary config and utils
FROM alpine:latest
WORKDIR /opt/app/

COPY  --from=BUILDER /go/src/github.com/opensourceways/xihe-server/xihe-server /opt/app

ENTRYPOINT ["/opt/app/xihe-server"]
