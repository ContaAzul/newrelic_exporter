FROM alpine:3.8 as alpine

RUN apk add -U --no-cache ca-certificates

FROM scratch
EXPOSE 9112
WORKDIR /
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY newrelic_exporter .
ENTRYPOINT ["./newrelic_exporter"]