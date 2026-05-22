FROM golang:1.25-alpine3.23 AS builder
RUN apk update && apk upgrade
RUN apk add git
COPY go.* /icanhazallips/
COPY ./vendor /icanhazallips/vendor
COPY *.go /icanhazallips/
RUN cd /icanhazallips && CGO_ENABLED=0 go build -mod=vendor -trimpath -ldflags="-s -w"

# FROM alpine:3.23
# COPY --from=builder /icanhazallips/icanhazallips /icanhazallips/icanhazallips
# ENTRYPOINT ["/icanhazallips/icanhazallips"]

FROM gcr.io/distroless/static-debian13:nonroot
WORKDIR /icanhazallips
COPY --from=builder /icanhazallips/icanhazallips /icanhazallips/icanhazallips
USER 65532:65532
ENTRYPOINT ["/icanhazallips/icanhazallips"]
