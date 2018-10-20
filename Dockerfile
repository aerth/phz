# build
FROM golang:1.11-alpine as builder

RUN apk add --no-cache make gcc musl-dev git

ADD . /phz
RUN cd /phz && make

# deploy
FROM alpine:latest

# get runtime deps
RUN apk add --no-cache ca-certificates curl

# copy binaries
COPY --from=builder /phz/phzd /usr/local/bin/
COPY --from=builder /phz/phz-cli /usr/local/bin/
COPY --from=builder /phz/_example/global.toml.default /etc/phz.toml

# serve it
EXPOSE 8080
CMD ["phzd", "-conf", "/etc/phz.toml"]
