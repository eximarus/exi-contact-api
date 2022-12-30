# exi-contact-api

## Environment

make sure the following environment variables are set
locally they are expected to be in a .env file in the project root

```
TARGET_EMAIL=
SMTP_USER=
SMTP_PASSWORD=
SMTP_HOST=
SMTP_PORT=
```

## Local setup

running on port 8080

```
go run ./cmd/local/main.go
```

or

```
docker-compose up
```

## Deployment

make sure the cdk is bootstrapped

```
cd cdk
cdk bootstrap
```

then you can run the deploy script from the repo root

```
./scripts/deploy.sh
```
