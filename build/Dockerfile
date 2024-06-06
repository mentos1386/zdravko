FROM gcr.io/distroless/static-debian12:latest as production

COPY zdravko /bin/zdravko
COPY LICENSE /LICENSE
COPY README.md /README.md

# Zdravko Server
ENV PORT=8080
EXPOSE 8080
# Temporal UI Server
EXPOSE 8223
# Temporal GRPC Server
EXPOSE 7233

# Volume to persist sqlite databases
#  as well as the keyvalue database.
ENV SQLITE_DATABASE_PATH=/data/zdravko.db
ENV TEMPORAL_DATABASE_PATH=/data/temporal.db
ENV KEYVALUE_DATABASE_PATH=/data/keyvalue.db
VOLUME /data

ENTRYPOINT ["/bin/zdravko"]
CMD ["--server", "--temporal", "--worker"]
