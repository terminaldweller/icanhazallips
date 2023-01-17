FROM alpine:3.17 as builder
ENV GOPROXY=https://goproxy.io
RUN apk update && apk upgrade
RUN apk add go git
ENV GOPROXY=https://goproxy.io
COPY go.* /icanhazallips/
RUN cd /icanhazallips && go mod download
COPY *.go /icanhazallips/
RUN cd /icanhazallips && go build

FROM alpine:3.17 as certbuilder
RUN apk add openssl
WORKDIR /certs
RUN openssl req -nodes -new -x509 -subj="/C=US/ST=Denial/L=springfield/O=Dis/CN=localhost" -keyout server.key -out server.cert

FROM gcr.io/distroless/static-debian11
# FROM alpine:3.17
COPY --from=certbuilder /certs /certs
COPY --from=builder /icanhazallips/icanhazallips /icanhazallips/icanhazallips
ENTRYPOINT ["/icanhazallips/icanhazallips"]
