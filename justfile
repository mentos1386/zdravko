# Always use devbox environment to run commands.
set shell := ["devbox", "run"]
# Load dotenv
set dotenv-load

# Load public and private keys
export JWT_PRIVATE_KEY := `cat jwt.private.pem`
export JWT_PUBLIC_KEY := `cat jwt.public.pem`

GIT_SHA := `git rev-parse --short HEAD`
DOCKER_IMAGE := "ghcr.io/mentos1386/zdravko:sha-"+GIT_SHA
STATIC_DIR := "./web/static"

# Build the application
build:
  docker build -f build/Dockerfile -t {{DOCKER_IMAGE}} .

# Run Docker application.
run-docker:
  docker run -p 8080:8080 \
  -e  SESSION_SECRET \
  -e OAUTH2_CLIENT_ID \
  -e OAUTH2_CLIENT_SECRET \
  -e OAUTH2_ENDPOINT_TOKEN_URL \
  -e OAUTH2_ENDPOINT_AUTH_URL \
  -e OAUTH2_ENDPOINT_USER_INFO_URL \
  -e OAUTH2_ENDPOINT_LOGOUT_URL \
  {{DOCKER_IMAGE}}

# Run full development environment
run:
  devbox services up

run-worker:
  go build -o dist/zdravko cmd/zdravko/main.go
  ./dist/zdravko --worker=true --server=false --temporal=false

# Generates new jwt key pair
generate-jwt-key:
  openssl genrsa -out jwt.private.pem 2048
  openssl rsa -pubout -in jwt.private.pem -out jwt.public.pem

# Start zdravko
run-zdravko:
  go build -o dist/zdravko cmd/zdravko/main.go
  ./dist/zdravko --worker=false

# Deploy the application to fly.io
deploy:
  fly deploy --ha=false -c deploy/fly.toml -i {{DOCKER_IMAGE}}

# Start devbox shell
shell:
  devbox shell

# Generate and download all external dependencies.
generate:
  rm -rf internal/models/query/*
  go generate ./...

_tailwindcss-build:
  tailwindcss build -c build/tailwind.config.js -i {{STATIC_DIR}}/css/main.css -o {{STATIC_DIR}}/css/tailwind.css

_htmx-download:
  mkdir -p  {{STATIC_DIR}}/js
  curl -sLo {{STATIC_DIR}}/js/htmx.min.js https://unpkg.com/htmx.org/dist/htmx.min.js

_feather-icons-download:
  mkdir -p {{STATIC_DIR}}/icons
  curl -sLo {{STATIC_DIR}}/icons/feather-sprite.svg https://unpkg.com/feather-icons/dist/feather-sprite.svg

_generate-gorm:
  go run tools/generate/main.go
