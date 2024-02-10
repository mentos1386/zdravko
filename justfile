# Always use devbox environment to run commands.
set shell := ["devbox", "run"]

STATIC_DIR := "./internal/static"

# Start devbox shell
shell:
  devbox shell

tailwindcss-build:
  tailwindcss build -i {{STATIC_DIR}}/css/main.css -o {{STATIC_DIR}}/css/tailwind.css

htmx-download:
  mkdir -p  {{STATIC_DIR}}/js
  curl -sLo {{STATIC_DIR}}/js/htmx.min.js https://unpkg.com/htmx.org/dist/htmx.min.js

feather-icons-download:
  mkdir -p {{STATIC_DIR}}/icons
  curl -sLo {{STATIC_DIR}}/icons/feather-sprite.svg https://unpkg.com/feather-icons/dist/feather-sprite.svg

generate:
  go generate ./...

# Start temporal which is accassible at http://localhost:8233/
run-temporal:
 temporal server start-dev

# Start web server accessible at http://localhost:8080/
run-server:
  go build -o dist/server cmd/server/main.go
  ./dist/server

# Run worker
run-worker:
  go build -o dist/worker cmd/worker/main.go
  ./dist/worker

# Run full development environment
run:
  devbox services up
