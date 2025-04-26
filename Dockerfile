FROM golang:1.23

WORKDIR ${GOPATH}/lost-items-service/
COPY . ${GOPATH}/lost-items-service/

RUN go build -o /build ./cmd/lost-items-service/ \
    && go clean -cache -modcache

EXPOSE 8080

CMD ["/build"]