# mtw

mtw is the proxy that converts email notifications into webhooks.

## Quick start

### With Docker

```bash
docker run -e "SECRET=mysecret" -e "DOMAIN=localhost.lan" -v ./data:/db -p "8080:8080" -p "25:25" -d zenentonal/mtw:v0.0.5
```

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
