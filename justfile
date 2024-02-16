# Always use devbox environment to run commands.
set shell := ["devbox", "run"]
# Load dotenv
set dotenv-load

STATIC_DIR := "./web/static"

# Build the application
build:
  docker build -f build/Dockerfile -t ghcr.io/mentos1386/zdravko:latest .

# Run full development environment
run:
  devbox services up

# Start zdravko
run-zdravko:
  go build -o dist/zdravko cmd/zdravko/main.go
  ./dist/zdravko

# Deploy the application to fly.io
deploy:
  fly deploy -c deploy/fly.toml

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
