FROM alpine:3.21 as builder
ENV GOPROXY=https://goproxy.io
RUN apk update && apk upgrade
RUN apk add go git
ENV GOPROXY=https://goproxy.io
COPY go.* /icanhazallips/
RUN cd /icanhazallips && go mod download
COPY *.go /icanhazallips/
RUN cd /icanhazallips && go build

FROM alpine:3.21
COPY --from=builder /icanhazallips/icanhazallips /icanhazallips/icanhazallips
ENTRYPOINT ["/icanhazallips/icanhazallips"]
