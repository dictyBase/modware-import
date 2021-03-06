FROM golang:1.16.4-buster AS builder
LABEL maintainer="Siddhartha Basu <siddhartha-basu@northwestern.edu>"
LABEL maintainer="Eric Hartline <eric.hartline@northwestern.edu>"
ENV GO111MODULE=auto \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
RUN apt-get -qq update \
    && apt-get -yqq install upx git
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
RUN strip gmake \
    && upx -q -9 gmake
RUN strip importer \
    && upx -q -9 importer

FROM gcr.io/distroless/static
COPY --from=builder /go/modware-import/gmake /usr/local/bin/
COPY --from=builder /go/modware-import/importer /usr/local/bin/
