# Always use devbox environment to run commands.
set shell := ["devbox", "run"]
# Load dotenv
set dotenv-load

STATIC_DIR := "./web/static"

# Run full development environment
run:
  devbox services up

# Start temporal which is accassible at http://localhost:8233/
run-temporal:
  go build -o dist/temporal cmd/temporal/main.go
  ./dist/temporal

# Start web server accessible at http://localhost:8080/
run-server:
  go build -o dist/server cmd/server/main.go
  ./dist/server

# Run worker
run-worker:
  go build -o dist/worker cmd/worker/main.go
  ./dist/worker

# Deploy the application to fly.io
deploy:
  fly deploy

# Start devbox shell
shell:
  devbox shell

# Generate and download all external dependencies.
generate:
  go generate ./...

_tailwindcss-build:
  tailwindcss build -i {{STATIC_DIR}}/css/main.css -o {{STATIC_DIR}}/css/tailwind.css

_htmx-download:
  mkdir -p  {{STATIC_DIR}}/js
  curl -sLo {{STATIC_DIR}}/js/htmx.min.js https://unpkg.com/htmx.org/dist/htmx.min.js

_feather-icons-download:
  mkdir -p {{STATIC_DIR}}/icons
  curl -sLo {{STATIC_DIR}}/icons/feather-sprite.svg https://unpkg.com/feather-icons/dist/feather-sprite.svg

_generate-gorm:
  go run tools/generate/main.go
