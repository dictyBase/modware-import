FROM golang:1.21.6-bullseye AS builder
LABEL maintainer="Siddhartha Basu <siddhartha-basu@northwestern.edu>"
ENV CGO_ENABLED=0 \
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
RUN go build  \
	-a \
	-ldflags "-s -w -extldflags '-static'" \
	-installsuffix cgo \
	-tags netgo \
	-o importer \ 
	cmd/import/main.go
RUN upx -q -9 gmake \
	&& upx -q -9 importer

RUN go build  \
	-a \
	-ldflags "-s -w -extldflags '-static'" \
	-installsuffix cgo \
	-tags netgo \
	-o baserow \ 
	cmd/baserow/main.go
RUN upx -q -9 gmake \
	&& upx -q -9 baserow

RUN go build  \
	-a \
	-ldflags "-s -w -extldflags '-static'" \
	-installsuffix cgo \
	-tags netgo \
	-o content \ 
	cmd/loader/main.go
RUN upx -q -9 gmake \
	&& upx -q -9 content

FROM gcr.io/distroless/static
COPY --from=builder /go/modware-import/gmake /usr/local/bin/
COPY --from=builder /go/modware-import/importer /usr/local/bin/
COPY --from=builder /go/modware-import/baserow /usr/local/bin/
COPY --from=builder /go/modware-import/content /usr/local/bin/
