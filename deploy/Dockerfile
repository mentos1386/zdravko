# syntax=docker/dockerfile:1
ARG ZDRAVKO_VERSION=main

FROM ghcr.io/mentos1386/zdravko:${ZDRAVKO_VERSION}
RUN apt-get update -y \
  && apt-get install -y ca-certificates fuse3 sqlite3

COPY --from=flyio/litefs:0.5 /usr/local/bin/litefs /usr/local/bin/litefs

COPY deploy/litefs.yaml /etc/litefs.yml
COPY deploy/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
