# letsgo

A tool to help build go web applications quickly.

Running this will create a go app with a Postgres database for local development, all in Docker containers.

## Usage
### Prerequisites
`letsgo` requires Go 1.21 or higher installed. Your app will run Go 1.22, but this will be containerized so it does not need to be installed locally.
`letsgo` also assumes that `docker compose` will run with a reasonably up to date version (using 3.8 of the compose file API as a minimum).


### Building
After cloning the repository, run `make build` to build, then move `letsgo` to somewhere in your path.

### Using
```letsgo make $projectname $projectRepository```

Example:
```letsgo make something-cool github.com/nmusey/something-cool```

Everything is Dockerized so you shouldn't have to install any dependencies locally.

### Migrations
`letsgo` uses [migrate](https://github.com/golang-migrate/migrate) for database migrations. By default these are run every time you restart the server. This is convenient for local development but should be migrated when moving to production.

`migrate` is not necessary locally but is preferred. You can always make migrations by copying the default files, but installing and running the tool is easier.

### Templating
`letsgo` uses [templ](https://templ.guide/) for templating. Hot reloading with Air is enabled by default.


## Future Features
- Support for API versioning is planned in the near future
- Support for a caching database is planned
- Support for other databases is not planned right now
- Support for other frameworks is not planned right now
