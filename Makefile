COVER_FILE=$(CURDIR)/coverage.out
BIN_DIR=$(CURDIR)/bin

PKG="github.com/a-novel/auth-service"

PKG_LIST=$(shell go list $(PKG)/... | grep -v /vendor/)

# Runs the test suite.
test:
	POSTGRES_URL=$(POSTGRES_URL_TEST) ENV="test" \
		gotestsum --packages="./..." --junitfile report.xml --format pkgname -- -count=1 -p 1 -v -coverpkg=./...

# Runs the test suite in race mode.
race:
	POSTGRES_URL=$(POSTGRES_URL_TEST) ENV="test" \
		gotestsum --packages="./..." --format pkgname -- -race -count=1 -p 1 -v -coverpkg=./...

# Run the test suite in memory-sanitizing mode. This mode only works on some Linux instances, so it is only suitable
# for CI environment.
msan:
	POSTGRES_URL=$(POSTGRES_URL_TEST) ENV="test" \
		env CC=clang env CXX=clang++ gotestsum --packages="./..." --format testname -- -msan -short $(PKG_LIST) -p 1

db-setup:
	psql -h localhost -p 5432 -U postgres agora -a -f init.sql

# Plugs into the development database.
db:
	psql -h localhost -p 5432 -U users agora_users

# Plugs into the test database.
db-test:
	psql -h localhost -p 5432 -U test agora_users_test

run:
	go run ./cmd/api/main.go

run-internal:
	go run ./cmd/api-internal/main.go

rotate-keys:
	curl -X POST http://localhost:20040/cloud/rotate-keys

.PHONY: all test race msan db db-test run run-internal
