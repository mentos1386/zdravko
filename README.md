# `zdravko`

Golang selfhosted Status/Healthcheck monitoring app.

Mostly just a project to test [temporal.io](https://temporal.io/).

### Roadmap
 - [x] SSO Support for authentication.
 - [x] SQLite for database.
    - This means for main app db as well as temporal db.
 - [x] Single binary.
    - One binary to run worker, server and temporal all together.
 - [x] Abbility for multiple workers.
    - Spread workers across regions to monitor latency from different locations.
 - [x] Use [k6](https://github.com/grafana/k6) for checks, so that they can be written in javascript.
 - [ ] History and working home page.
 - [ ] Edit/Delete operations for healthchecks and workers.
 - [ ] CronJob Healthchecks (via webhooks).
 - [ ] Notifications (webhooks, slack, etc).

![Screenshot](docs/screenshot.png)
Demo is available at https://zdravko.mnts.dev.

# Development

### Dependencies
 * [devbox](https://www.jetpack.io/devbox)
 * [justfile](https://github.com/casey/just) (optional, `devbox run -- just` can be used instead)

```sh
# Configure
cp example.env .env

# Generate JWT key
just generate-jwt-key

# Start development environment
just run
```

### License
Under AGPL, see [LICENSE](LICENSE) file.
