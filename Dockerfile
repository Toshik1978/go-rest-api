## Builder
FROM golang:1.13-alpine3.10 AS build

RUN apk --no-cache add make git

ENV GOPATH /.go/
ENV GOPROXY https://proxy.golang.org/,direct

# Go modules
WORKDIR /app
COPY go.mod /app
COPY go.sum /app
COPY Makefile /app
RUN make modules

# Build app
COPY . /app
RUN make build

## Destination image
FROM alpine:3.10

RUN apk --no-cache add ca-certificates supervisor tzdata

COPY --from=build /app/go-rest-api /opt/go-rest-api/go-rest-api
COPY --from=build /app/configs/go-rest-api.conf.yaml /opt/go-rest-api/configs/go-rest-api.conf.yaml
COPY .docker/supervisord.conf /etc/supervisord.conf

VOLUME ["/var/log/supervisor", "/opt/go-rest-api/configs"]
WORKDIR /opt/go-rest-api/
EXPOSE 8080

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor.conf"]
