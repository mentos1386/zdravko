# Always use devbox environment to run commands.
set shell := ["devbox", "run"]
# Load dotenv
set dotenv-load

# Load public and private keys
export JWT_PRIVATE_KEY := `cat jwt.private.pem || echo ""`
export JWT_PUBLIC_KEY := `cat jwt.public.pem || echo ""`

GIT_SHA := `git rev-parse --short HEAD`
DOCKER_IMAGE := "ghcr.io/mentos1386/zdravko:sha-"+GIT_SHA

STATIC_DIR := "./web/static"

OS := "linux darwin"
ARCH := "amd64 arm64"

export CGO_ENABLED := "0"
import 'build/Justfile'

_default:
  @just --list

# Run full development environment
run:
  watchexec -r -e tmpl,css just tailwindcss | sed -e 's/^/tailwind: /;' &
  sleep 1
  just run-temporal | sed -e 's/^/temporal: /;' &
  sleep 1
  watchexec -r -e go,tmpl,css just run-server

# Start worker
run-worker:
  go run cmd/zdravko/main.go --worker

# Start server
run-server:
  go run cmd/zdravko/main.go --server

# Start temporal
run-temporal:
  go run cmd/zdravko/main.go --temporal

# Test
test:
  go test -v ./...

# Generates new jwt key pair
generate-jwt-key:
  openssl genrsa -out jwt.private.pem 2048
  openssl rsa -pubout -in jwt.private.pem -out jwt.public.pem

# Run Docker application.
run-docker: build-docker
  docker run -p 8080:8080 \
  -it --rm \
  -e SESSION_SECRET \
  -e OAUTH2_CLIENT_ID \
  -e OAUTH2_CLIENT_SECRET \
  -e OAUTH2_ENDPOINT_TOKEN_URL \
  -e OAUTH2_ENDPOINT_AUTH_URL \
  -e OAUTH2_ENDPOINT_USER_INFO_URL \
  -e OAUTH2_ENDPOINT_LOGOUT_URL \
  -e JWT_PRIVATE_KEY \
  -e JWT_PUBLIC_KEY \
  -e WORKER_GROUP_TOKEN \
  {{DOCKER_IMAGE}} --server --temporal --worker

# Start Sqlite web client
sqlite-web:
  sqlite_web zdravko.db

# New migration file
migration-new name:
  #!/bin/bash
  FILENAME="database/sqlite/migrations/`date --iso-8601`_{{name}}.sql"

  cat <<EOF > $FILENAME
  -- +migrate Up
  -- SQL in section 'Up' is executed when this migration is applied

  -- +migrate Down
  -- SQL in section 'Down' is executed when this migration is rolled back
  EOF

  echo "Created migration file: $FILENAME"

update-dependencies:
  go get -u -t \
    go.k6.io/k6@v0.51.0 \
    github.com/temporalio/ui-server/v2@v2.23.0 \
    go.temporal.io/server@v1.23.0 \
    go.temporal.io/api@v1.29.2 \
    go.temporal.io/sdk@v1.26.0 \
    ./...
  go mod tidy

# Run go generate and process tailwindcss
generate: tailwindcss
  go generate ./...

tailwindcss:
  mkdir -p {{STATIC_DIR}}/css
  tailwindcss build -c build/tailwind.config.js -i {{STATIC_DIR}}/css/main.css -o {{STATIC_DIR}}/css/tailwind.css

static-dependencies:
  npm install

static-clean:
  find {{STATIC_DIR}} -type f -not -path '{{STATIC_DIR}}/static.go' -not -path '{{STATIC_DIR}}/css/*' -exec rm -f {} \;

static: static-dependencies static-clean tailwindcss
  # HTMX
  mkdir -p {{STATIC_DIR}}/js
  cp node_modules/htmx.org/dist/htmx.min.js {{STATIC_DIR}}/js/htmx.min.js

  # Monaco
  cp -r node_modules/monaco-editor/min/* {{STATIC_DIR}}/monaco
  # We only care about javascript language
  find {{STATIC_DIR}}/monaco/vs/basic-languages/ \
    -type d \
    -not -name 'javascript' \
    -not -name 'typescript' \
    -not -name 'yaml' \
    -not -name 'basic-languages' \
    -prune -exec rm -rf {} \;

  # Feather Icons
  mkdir -p {{STATIC_DIR}}/icons
  cp node_modules/feather-icons/dist/feather-sprite.svg {{STATIC_DIR}}/icons/feather-sprite.svg
