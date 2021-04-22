FROM golang:1.16.3-alpine3.13
LABEL maintainer="Siddhartha Basu <siddhartha-basu@northwestern.edu>"
ENV GOPROXY https://proxy.golang.org
RUN apk add --no-cache git build-base
# install mage
RUN git clone https://github.com/magefile/mage \
    && cd mage \
    && go run bootstrap.go
RUN mkdir -p /modware-import
WORKDIR /modware-import
COPY go.mod ./
COPY go.sum ./
COPY magefile.go ./
RUN go mod download
ADD cmd cmd
ADD internal internal
RUN go build -o app cmd/import/main.go
RUN mage -compile ./gmake


FROM alpine:3.13.2
RUN apk --no-cache add ca-certificates tzdata \
    && cp /usr/share/zoneinfo/America/Chicago /etc/localtime \
    && echo "America/Chicago" > /etc/timezone
COPY --from=0 /modware-import/app /usr/local/bin/
COPY build/package/dsc-runner.sh /usr/local/bin/dsc-runner
RUN chmod +x /usr/local/bin/dsc-runner \
    && mkdir -p /build \
    && cd /build
COPY --from=0 /modware-import/app importer
COPY --from=0 /modware-import/gmake .
RUN chmod +x importer gmake
