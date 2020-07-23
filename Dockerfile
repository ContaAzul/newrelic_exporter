FROM golang:1.10-alpine as build
WORKDIR /go/src/github.com/ContaAzul/newrelic_exporter
COPY . /go/src/github.com/ContaAzul/newrelic_exporter
RUN apk --no-cache add ca-certificates curl git
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure -v -vendor-only
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o newrelic_exporter

FROM scratch
EXPOSE 9112
WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/github.com/ContaAzul/newrelic_exporter/newrelic_exporter /

ENTRYPOINT ["/newrelic_exporter"]
