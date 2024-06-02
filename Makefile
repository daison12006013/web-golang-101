run:
	LOG_LEVEL="debug" \
	JWT_KEY="YourJWTKeyHere" \
	APP_KEY="YourAPPKeyHere" \
	MAILGUN_DOMAIN="mailer.acme.com" \
	MAILGUN_API_KEY="YourMailgunAPIKey" \
	MAIL_FROM="no-reply@acme.com" \
	DB_STRING="postgres://postgres:YourPostgresPassword@localhost:5432/acme_dev?sslmode=disable" \
	go run ./cmd/http

build-individual:
	(cd ./cmd/$(dir)/ && GOOS=darwin GOARCH=amd64 go build -o main-darwin.bin)
	(cd ./cmd/$(dir)/ && GOOS=linux GOARCH=arm64 go build -o main-linux-arm64.bin)

build:
	for dir in $$(ls -d ./cmd/*/); do \
		make build-individual dir=$$(basename $$dir); \
	done

clean:
	rm -f ./cmd/*/*.bin

deploy:
	fly deploy

remove-ds-store:
	find . -name .DS_Store -delete

swagger:
	swag init -d cmd/http/ --parseDependency --parseInternal

dbproxy:
	make dbstart && fly proxy 5432 -a pg-prod-acme-com

dbstart:
	(cd migrations && fly machines start)

dbstop:
	(cd migrations && fly machines stop)

goose:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="postgres://postgres:YourPostgresPassword@localhost:5432/acme_dev" goose -dir=migrations $(filter-out $@,$(MAKECMDGOALS))

ci:
	go mod tidy
	go mod verify

test:
	go test -v ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

%:
	@:

