FROM golang:1.17.6 as BUILDER

WORKDIR /build
ADD . /build/
RUN go mod tidy && make

FROM alpine:3.15

RUN apk add ca-certificates

COPY --from=BUILDER /build/bin/itsnotgolang /usr/local/bin

ENTRYPOINT ["itsnotgolang"]
