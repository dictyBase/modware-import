FROM golang:1.11.10-alpine3.9
LABEL maintainer="Siddhartha Basu <siddhartha-basu@northwestern.edu>"
RUN apk add --no-cache git build-base
RUN mkdir -p /modware-import
WORKDIR /modware-import
ADD cmd cmd
ADD internal internal
COPY go.mod go.sum ./
RUN go build -o app cmd/import/main.go

FROM alpine:3.9
RUN apk --no-cache add ca-certificates
COPY --from=0 /modware-import/app /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/app"]
