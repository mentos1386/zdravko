# syntax=docker/dockerfile:1
ARG GO_VERSION=1.21
FROM golang:${GO_VERSION} as builder
# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . ./

# Build
RUN CGO_ENABLED=1 GOOS=linux go build -o /bin/zdravko cmd/zdravko/main.go

###
# Final production
FROM gcr.io/distroless/base-debian12:nonroot as production
COPY --from=builder /bin/zdravko /bin/zdravko
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
VOLUME /data

ENV DATABASE_PATH=/data/zdravko.db
ENV TEMPORAL_DATABASE_PATH=/data/temporal.db

ENTRYPOINT ["/bin/zdravko"]
