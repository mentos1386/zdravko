# `zdravko`

Golang selfhosted Status/Healthcheck monitoring app.

Mostly just a project to test [temporal.io](https://temporal.io/).

### Expected ~features~ things
 - [ ] SSO or just expect to be run behind a auth proxy.
 - [ ] Abbility for multiple instances/workers.
   - Otherwise using a cronjob would do the job, no need for temporal.
 - [ ] Some nice UI to try out [htmx](https://htmx.org/).

![Screenshot](docs/screenshot.png)
Demo is available at https://zdravko.fly.io.

# Development

### Dependencies
 * [devbox](https://www.jetpack.io/devbox)
 * [justfile](https://github.com/casey/just)

```sh
# Start development environment
just run
```
