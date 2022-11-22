FROM golang:1.19-alpine AS builder

RUN apk add gcc g++ musl-dev

WORKDIR /tmp/art

ADD . .

RUN go build -o art .

FROM alpine:3.16

RUN apk add libc6-compat

COPY --from=builder /tmp/art /art

RUN chmod a+x /art

RUN mkdir -p out/

EXPOSE 3000

ENTRYPOINT ["/art"]
