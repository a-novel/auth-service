# Auth service

Manage platform users and authentication.

## Prerequisites

 - Download [Go](https://go.dev/doc/install)
 - Install [Mockery](https://vektra.github.io/mockery/latest/installation/)
 - Clone [go-framework](https://github.com/a-novel/go-framework)
   - From the framework, run `docker compose up -d`

## Installation

Create a env file.

> Ask an admin for the Sendgrid API key.

```bash
touch .envrc
```
```bash
printf 'export POSTGRES_URL="postgres://users@localhost:5432/agora_users?sslmode=disable"
export POSTGRES_URL_TEST="postgres://test@localhost:5432/agora_users_test?sslmode=disable"
export SENDGRID_API_KEY="xxxxxxxxx"
' > .envrc
```
```bash
direnv allow .
```

Set the database up.
```bash
make db-setup
```

Run the internal API.
```bash
make run-internal
```

> Check the API is up by running `curl http://localhost:20040/ping`.

In a new terminal, create local keys.
```bash
make rotate-keys
```
> You can now kill the internal API. You may repeat these 2 last steps every time you want to rotate your local
> keys.

## Commands

### Run the API

```bash
make run
```
```bash
curl http://localhost:2040/ping
# Or curl http://localhost:2040/healthcheck
```

### Run the internal API

```bash
make run-internal
```
```bash
curl http://localhost:20040/ping
# Or curl http://localhost:20040/healthcheck
```

### Rotate local keys

```bash
make run-internal
```
In another terminal.
```bash
make rotate-keys
```

### Run tests

```bash
make test
```

### Update mocks

```bash
mockery
```

### Open a postgres console

```bash
make db
# Or make db-test
```
