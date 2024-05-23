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
 - [x] History and working home page.
    - Kinda working atm. ~But look if all the data could be stored/fetched from temporal.~
 - [x] Edit/Delete operations for healthchecks and workers.
 - [ ] CronJob Healthchecks (via webhooks).
 - [ ] Notifications (webhooks, slack, etc).
 - [ ] Incidents (based on script that is triggered by monitors/crobjobs).
 - [ ] Prepare i18n.
 - [ ] Alpha Version (1H 2024)
 - [ ] ??
 - [ ] Beta Version (2H 2024)
 - [ ] ??
 - [ ] Stable Release (2025)

![Screenshot](docs/screenshot.png)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fmentos1386%2Fzdravko.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fmentos1386%2Fzdravko?ref=badge_shield)
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


[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fmentos1386%2Fzdravko.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fmentos1386%2Fzdravko?ref=badge_large)