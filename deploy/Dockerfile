# syntax=docker/dockerfile:1
ARG ALPINE_VERSION=3.14

FROM alpine:${ALPINE_VERSION} as temporal
RUN apk add ca-certificates fuse3 sqlite
COPY --from=flyio/litefs:0.5 /usr/local/bin/litefs /usr/local/bin/litefs
COPY litefs.yaml /etc/litefs.yaml

COPY --from=ghcr.io/mentos1386/zdravko /zdravko /zdravko

ENTRYPOINT litefs mount
