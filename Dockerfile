FROM golang:1.25-alpine3.23 AS builder
ENV GOPROXY=https://goproxy.io
RUN apk update && apk upgrade
RUN apk add go git
COPY go.* /icanhazallips/
RUN cd /icanhazallips && go mod download
COPY *.go /icanhazallips/
RUN cd /icanhazallips && go build

FROM alpine:3.23
COPY --from=builder /icanhazallips/icanhazallips /icanhazallips/icanhazallips
ENTRYPOINT ["/icanhazallips/icanhazallips"]
