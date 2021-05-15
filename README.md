# URL Shortener

A simple url shortener using Base62 conversion.

## Requirements

- Go **1.16+**
- docker-compose **1.27.0+**

## Features

- Works without JavaScript ☑️
- Dark mode 🌒
- Single binary (`go:embed`)

## Usage

```sh
$ ./bin/go-url-shortener -help
Usage of ./bin/go-url-shortener:
  -port string
        server port (default "4000")
  -public string
        public url prefix (default "http://localhost:4000/")
```

## Running the application

```sh
# Create .env file
$ make -B .env env=local
# Download go modules
$ go mod download
# Start linked services (postgres, etc.)
$ make deps
# Apply database migrations
$ ./scripts/migrate.sh up

# Run in watch mode (restart on file change)
$ make dev
# Or run in normal mode
$ go run ./cmd/go-url-shortener/
```

## Testing

```sh
$ make test
```

## Dockerizing

```sh
$ make -B .env env=docker
$ ./scripts/migrate.sh up
$ docker-compose up app
```

## Database migrations

[migrate](https://github.com/golang-migrate/migrate) - CLI and Golang library.

### Up

```sh
$ ./scripts/migrate.sh up
```

### Down

```sh
$ ./scripts/migrate.sh down
```

### Create

```sh
$ ./scripts/migrate.sh create -dir /migrations -ext sql -seq create_new_table

Output:
migrations/XXXXXX_create_new_table.up.sql
migrations/XXXXXX_create_new_table.down.sql
```

## License

[MIT](/LICENSE)
