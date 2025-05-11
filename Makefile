PROJECT_NAME=mcp-memory

migration: export BUN_TIMESTAMP=$(shell date +%Y%m%d%H%M%S | head -c 14)
migration:
	touch internal/db/bun/migrations/${BUN_TIMESTAMP}_new.go
	cat internal/db/bun/migrations/migration.go.tmpl > internal/db/bun/migrations/${BUN_TIMESTAMP}_new.go
	sed -i 's/CHANGEME/${BUN_TIMESTAMP}/g' internal/db/bun/migrations/${BUN_TIMESTAMP}_new.go

snapshot:
	goreleaser build --clean --snapshot

test: tidy fmt
	go test -cover ./...

tidy:
	go mod tidy

fmt:
	@echo formatting
	@go fmt $(shell go list ./... | grep -v /vendor/)
