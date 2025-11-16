include .env
export

MIGRATIONS_DIR=./migrations
DSN = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)
MIGRATE = migrate -path $(MIGRATIONS_DIR) -database "$(DSN)" 
.PHONY: migrate-up migrate-down  migrate-version  migrate-down-n migrate-goto migrate-dsn
migrate-up:
	$(MIGRATE) up
migrate-down:
	$(MIGRATE) down
migrate-down-n:
	$(MIGRATE) down $(N) 
migrate-goto:
	$(MIGRATE) goto $(V)
migrate-version:
	$(MIGRATE) version
migrate-dsn:
	@echo $(DSN)

migrate-reset:
	$(MIGRATE) drop -f
	$(MIGRATE) up
vet:
	go vet ./...
run:
	go run ./cmd/server