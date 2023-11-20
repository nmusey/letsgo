# letsgo

A tool to help build go web applications quickly.

Running this will create a Fiber app with a Postgres database for local development, all in Docker containers.

## Usage
### Building
`letsgo` requires Go 1.21 or higher installed.
After cloning the repository, run `make build` to build, then move `letsgo` to somewhere in your path.

### Using
```letsgo make $projectname $projectRepository```

Example:
```letsgo make something-cool github.com/nmusey/something-cool```

Everything is Dockerized so you shouldn't have to install any dependencies, that's taken care of by the container.

### Migrations
`letsgo` uses [migrate](https://github.com/golang-migrate/migrate) for database migrations. By default these are run every time you restart the server.


## Future Features
- Support for caching databases is planned
- Support for other databases is possible
- Support for other frameworks is not planned right now
