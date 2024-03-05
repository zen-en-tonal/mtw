# mtw

mtw is the proxy that converts email notifications into webhooks.

## Quick start

### With Docker

1. Create a `.env` file with the following content
```bash
DOMAIN="localhost.lan"
DB_CONN="postgres://postgres:postgres@db:5432/postgres?sslmode=disable"
SECRET="mysecret"
```
2. Create a `docker-compose.yml` file with the following conetnt
```yml
version: '3'
services:
  mtw:
    image: zenentonal/mtw:v0.0.4
    restart: unless-stopped
    env_file: .env
    depends_on:
      db:
        condition: service_healthy
    ports:
      - "8080:8080"
      - "25:25"
  db:
    image: postgres
    restart: unless-stopped
    environment:
      POSTGRES_PASSWORD: postgres
    volumes:
      - ./db:/var/lib/postgresql/data
    healthcheck:
      test: 'pg_isready -U "${POSTGRES_USER:-postgres}"'
      interval: 10s
      timeout: 10s
      retries: 3
      start_period: 30s
```
3. Run `docker compose up -d`

## Tutorial

### Step 1. Make an address

```bash
curl -XPOST localhost:8080/address/user/alice \
     -H 'Authorization: Bearer mysecret'

{"address":"alice@localhost.lan"}
```

### Step 2. Make a webhook

```bash
curl -XPOST localhost:8080/webhook \
     -H 'Authorization: Bearer mysecret' \
     -H 'Content-Type: application/json' \
     --data-raw '
{
    "endpoint": "https://hooks.slack.com/services/xxxx",
    "method": "POST",
    "schema": "{\"text\": \"{{Escape .Text | Limit 3000}}\"}",
    "content_type": "application/json"
}'

{"id":"19116242-dfdc-4b94-bce6-0b4cc90ec372"}
```

### Step 3. Link an address to a webhook

```bash
curl -XPOST localhost:8080/address/alice@localhost.lan/webhook/19116242-dfdc-4b94-bce6-0b4cc90ec372 \
     -H 'Authorization: Bearer mysecret'
```

## Licence

MIT
