# letsgo

A tool to help build go web applications quickly.

Running this will create a Fiber app with a Postgres database for local development, all in Docker containers.

## Usage
```letsgo make $projectname $projectRepository```

Example:
```letsgo make something-cool github.com/nmusey/something-cool```

Everything is Dockerized so you shouldn't have to install any dependencies, that's taken care of by the container.

## Future Features
- Support for caching databases is planned
- Support for other databases is possible
- Support for other frameworks is not planned right now
- Support for a better database access layer is planned
