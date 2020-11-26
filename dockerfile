FROM library/golang:1.15-alpine
RUN apk update && apk --no-cache add bash git make ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

ARG CGO_ENABLED=0
ENV APP_NAME rosseti-innovation-back
ENV APP_DIR $GOPATH/src/github.com/sqsinformatique/$APP_NAME
WORKDIR $APP_DIR

COPY . .
RUN make build-stable && \
    cd /go/bin/; ln -s /go/bin/$APP_NAME /go/bin/service && \
    ln -s $APP_DIR/internal/db/migrations /go/bin/migrations

FROM scratch

COPY --from=0 /go/bin/$APP_NAME /usr/bin/$APP_NAME
COPY --from=0 /go/bin/service /usr/bin/service
COPY --from=0 /go/bin/migrations /usr/bin/migrations
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER 9000

# Set the entrypoint
ENTRYPOINT ["service"]