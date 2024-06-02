**Attention:** ***Feel free to experiment with this, but please be aware it's still under enhancement.***


# Web Golang 101

Web Golang 101 is an open-source project designed to provide a practical and direct approach to integrating Golang into your web platform. It offers a robust set of features including Sentry DSN, Zero Log, Mailgun, Goose Migration, SQLc, and AES-CBC-256 Encryption, all designed to align with the standards of popular web frameworks like Laravel.

The project now includes prebuilt endpoints for registration, email verification, login, and token refresh. These endpoints utilize JSON Web Tokens (JWT) for secure authentication and session management, serving as a comprehensive starting point for building secure and efficient web applications with Golang. This JWT-based approach ensures a stateless, scalable solution that can easily integrate with various front-end frameworks.

The project also integrates with LOG_SENTRY_LEVEL (default: "warn") and includes an encrypted users table where the email is stored as an encrypted value.

## Features

- :white_check_mark: Sentry DSN
- :white_check_mark: Zero Log
    - integrated with `LOG_SENTRY_LEVEL` (default: "warn")
- :white_check_mark: Mailgun
- :white_check_mark: Goose Migration
- :white_check_mark: SQLc
- :white_check_mark: AES-CBC-256 Encryption
    - Follows most of web frameworks standards such as Laravel
- :white_check_mark: Encrypted `users` table
    - This is where the `email` and password are stored in encrypted value
    - The `email_hash` is also stored as sha256 for filtering purposes

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
JWT_REFRESH_EXP="168d"
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

## Prebuilt Endpoints

**authorization:**

- :white_check_mark: (POS) /register
- :white_check_mark: (GET) /verify-email/{token}
- :white_check_mark: (POS) /login
- :white_check_mark: (POS) /refresh-token

**api key generator for your business:**

- :white_check_mark: (GET) /api-keys
- :white_check_mark: (POS) /api-keys
- :white_check_mark: (DEL) /api-keys/{key}

## Makefile

To run a local server

```bash
make run
```

To tunnel the database (via fly.io)

```bash
make dbproxy
```

Calling goose

```bash
make goose
```

## Cloud Services

- :white_check_mark: Fly.io
- :white_check_mark: AWS Fargate

## License

This project is licensed under the terms of the Apache 2.0 license. For more details, see the [LICENSE](LICENSE) file in the project root.
