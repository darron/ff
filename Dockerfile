FROM golang:1.21 as build
RUN mkdir /ff
RUN go install github.com/kevinburke/go-bindata/v4/...@latest
WORKDIR /ff
ADD . .
RUN make linux
RUN chmod a+x /ff/bin/ff

FROM cgr.dev/chainguard/wolfi-base:latest
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
LABEL org.label-schema.vcs-url="https://github.com/darron/ff"
RUN apk add --update --no-cache \
  ca-certificates && \
  rm -vf /var/cache/apk/*
WORKDIR /
COPY --from=build /ff/views /views
COPY --from=build /ff/public /public
COPY --from=build /ff/bin/ff /
ENTRYPOINT ["./ff", "service"]