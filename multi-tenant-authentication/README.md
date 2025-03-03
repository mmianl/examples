# Multi-Tenant Authentication with Zitadel

## Structure

This repository contains the following directories:

- `api`: The user and admin API for the Multi-Tenant Authentication service: [README](./api/README.md)

## Setup

Run the docker compose file to start the postgres instance and Zitadel.

```bash
docker compose up
```

Connect to the zitadel instance at [http://localhost:8080/](http://localhost:8080/) and login with the credentials `zitadel-admin@zitadel.localhost` and `Password1!`.
