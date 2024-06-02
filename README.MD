**Attention:** ***Feel free to experiment with this, but please be aware it's still under enhancement.***

[![Go Report Card](https://goreportcard.com/badge/github.com/daison12006013/web-golang-101)](https://goreportcard.com/report/github.com/daison12006013/web-golang-101) [![Go Reference](https://pkg.go.dev/badge/github.com/daison12006013/web-golang-101.svg)](https://pkg.go.dev/github.com/daison12006013/web-golang-101) [![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# Web Golang 101: A Comprehensive Framework for Golang Web Development

Web Golang 101 is an all-inclusive, robust framework designed to facilitate web development using the Go programming language. This project integrates a wide array of tools and features, providing a streamlined development process and a high degree of efficiency and security.

## Features

- [Sentry](https://sentry.io/)
- [Mailgun](https://app.mailgun.com/mg/dashboard)
- [Goose Migration](https://github.com/pressly/goose)
- [SQLc](https://github.com/sqlc-dev/sqlc)
- [Swagger API Doc](https://github.com/swaggo/swag)
- [AES-256-CBC](https://docs.anchormydata.com/docs/what-is-aes-256-cbc) Encryption
    - Follows most of web frameworks standards such as [Laravel](https://laravel.com)
    - Secured `users` table
        - column `email` is encrypted by default using AES-256-CBC
        - column `email_hash` is hashed using sha256
        - column `password` uses bcrypt
- [Zero Log](https://github.com/rs/zerolog)
    - integrated with sentry thru hooks `LOG_SENTRY_LEVEL` (default: "warn")

## ENV Vars

```
# app
APP_ENV="production"
APP_KEY="YourAPPKeyHere"
APP_DEBUG=true
APP_PORT="8080"

# jwt
JWT_ISSUER="web-golang-101"
JWT_EXP="1h"
JWT_REFRESH_EXP="168h"
JWT_KEY="YourJWTKeyHere"

# database
DB_STRING="postgres://postgres:YourPostgresPassword@localhost:5432/your_database_dev?sslmode=disable"
DB_MAX_OPEN_CONNS=20

# mailgun
MAIL_FROM="no-reply@acme.com"
MAILGUN_API_KEY="YourMailgunAPIKey"
MAILGUN_DOMAIN="mailer.acme.com"

# sentry
SENTRY_DSN=""
SENTRY_LEVEL="warn"

# log
LOG_APP_LEVEL="warn"
LOG_SENTRY_LEVEL="warn"
LOG_REQUEST_ENABLED=true

# web application firewall
WAF_ENABLED=true
WAF_RATE_LIMIT="5"
WAF_REJECT_REQUESTS_EXCEPT="acme.com"

# get requestor's ip
IP_CLIENT_HEADER_KEY="CF-Connecting-IP"

# cross origin resource
CORS_ORIGIN="*"
CORS_METHODS="POST, GET, OPTIONS, PUT, DELETE"
CORS_HEADERS="Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-API-KEY"
```

## Default service

- Fly.io
    - fly.yaml
- AWS Fargate
    - copilot

## Prebuilt Endpoints

- **authorization:**
    - (POS) /register
    - (GET) /verify-email/{token}
    - (POS) /login
    - (POS) /refresh-token
- **api key generator for your business:**
    - (GET) /api-keys
    - (POS) /api-keys
    - (DEL) /api-keys/{key}

## Reporting Security Vulnerabilities

We prioritize security in our project. If you find a security vulnerability, please email us directly at daison12006013@gmail.com instead of publicizing it. Provide as much detail as possible about the vulnerability, including its nature, reproduction steps, and potential impacts. We appreciate your help in enhancing our application's security. All reports are thoroughly investigated and addressed promptly and responsibly to maintain our software's safety and integrity.

## License

This project is licensed under the terms of the Apache 2.0 license. For more details, see the [LICENSE](LICENSE) file in the project root.
