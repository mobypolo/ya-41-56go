# Gofermart

Cumulative loyalty system.

[![go report card](https://goreportcard.com/badge/github.com/mobypolo/ya-41-56go?style=flat-square)](https://goreportcard.com/report/github.com/mobypolo/ya-41-56go)
[![test status](https://github.com/mobypolo/ya-41-56go/workflows/gophermart/badge.svg?branch=main "test status")](https://github.com/mobypolo/ya-41-56go/actions)

## Getting Started

Dependencies:

* Go `1.23`
* Docker for PostgreSQL
* Linux or macOS platform

### Startup

#### Gophermart

To build a gophermart, run in the terminal:

```bash
make build_gophermart # in the root directory of the project
```

To startup the gophermart, run in the terminal:

```bash
./gophermart # in the root directory of the project
```

#### Accrual

To build a accrual, run in the terminal:

```bash
make accrual # in the root directory of the project
```

To startup the accrual, run in the terminal:

```bash
./accrual # in the root directory of the project
```

#### Sartup envs & flags

* `JWT_SECRET_KEY` – JWT secret key.
* `JWT_LIFETIME` – JWT lifetime. Default `1h`
* `DATABASE_URI` | `-d` – PostgresSQL DSN
* `RUN_ADDRESS` | `-a` – HTTP server address. Default `localhost:8080`
* `LOG_MODE` – Logging mode. Default `dev`. Values `info`, `dev`, `warn`, `error`
* `CORS_ORIGINS` – List of CORS-origins. Default `http://localhost:3000`
* `WORKERS_COUNT` – Count of workers. Default `5`
* `SHUTDOWN_TIMEOUT` – Shutdown timeout. Default `5s`

## Development

### Updating the template

To be able to receive updates for autotests and other parts of the template, run the command:

```bash
git remote add -m master template https://github.com/yandex-praktikum/go-musthave-group-diploma-tpl.git
```

To update the autotest code, run the command:

```bash
git fetch template && git checkout template/master .github
```

Then add the changes to your repository.
