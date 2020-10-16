# Estuary (WIP)

This is the Go server that will replace my NodeJS server for Estuary, a productivity app that I've been working on.

## Set Up

To set up, you'll need the following in a .env file:

```
SESSION_SECRET
TEST_PSQL_PW
CSRF_AUTH_KEY
ENV
MG_API_KEY
MG_DOMAIN
```

## Running

Currently, you can locally run this server with `go run ./cmd/web`