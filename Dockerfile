#
# Cloudmanic Backup Dockerfile
#

FROM golang:1.9.2-alpine3.7

MAINTAINER Spicer Matthews <spicer@cloudmanic.com>

# Build App
RUN apk --no-cache add --virtual build-dependencies go git \
  && go get github.com/cloudmanic/cloudmanic-backup \
  && apk del --purge build-dependencies

# Install Helpful Programs
RUN apk --no-cache add bash mysql-client

RUN adduser -D -u 1000 backup

USER backup

WORKDIR /go/bin

CMD ["cloudmanic-backup"]