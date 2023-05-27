FROM golang:1.20 as build
RUN mkdir /ff
WORKDIR /ff
ADD . .
RUN make linux
RUN chmod a+x /ff/bin/ff

FROM cgr.dev/chainguard/alpine-base:latest
LABEL org.label-schema.vcs-url="https://github.com/darron/ff"
RUN apk add --update --no-cache \
  ca-certificates && \
  rm -vf /var/cache/apk/*
WORKDIR /
COPY --from=build /ff/views /views
COPY --from=build /ff/public /public
COPY --from=build /ff/bin/ff .
ENTRYPOINT ["/ff service"]