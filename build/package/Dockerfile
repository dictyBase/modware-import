FROM golang:1.16.4-alpine3.13 AS builder
LABEL maintainer="Siddhartha Basu <siddhartha-basu@northwestern.edu>"
RUN apk add --no-cache git 
ENV GO111MODULE=auto
# install mage
RUN git clone https://github.com/magefile/mage \
    && cd mage \
    && go run bootstrap.go \
    && mkdir -p /go/modware-import 
WORKDIR /go/modware-import
COPY go.mod ./
COPY go.sum ./
COPY magefile.go ./
RUN go mod download
ADD cmd cmd
ADD internal internal
RUN mage -compile ./gmake 
RUN go build -o importer cmd/import/main.go

FROM alpine:3.13.5
RUN apk --no-cache add ca-certificates tzdata \
    && cp /usr/share/zoneinfo/America/Chicago /etc/localtime \
    && echo "America/Chicago" > /etc/timezone
COPY --from=builder /go/modware-import/gmake /usr/local/bin/
COPY --from=builder /go/modware-import/importer /usr/local/bin/
