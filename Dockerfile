FROM golang:1.20 as build
RUN mkdir /ff
WORKDIR /ff
ADD . .
RUN make linux
RUN chmod a+x /ff/bin/ff

FROM cgr.dev/chainguard/alpine-base:latest
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
LABEL org.label-schema.vcs-url="https://github.com/darron/ff"
RUN apk add --update --no-cache \
  ca-certificates && \
  rm -vf /var/cache/apk/*
COPY --from=build /ff/views /bin/views
COPY --from=build /ff/public /bin/public
COPY --from=build /ff/bin/ff /bin/
WORKDIR /
ENTRYPOINT ["ff", "service"]