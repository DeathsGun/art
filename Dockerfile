FROM golang:1.19-alpine AS builder

RUN apk add gcc g++ musl-dev

WORKDIR /tmp/art

ADD . .

RUN go build -o art .

FROM alpine:3.16

COPY --from=builder /tmp/art /art

RUN mkdir -p out/

EXPOSE 3000

ENTRYPOINT ["/art"]
