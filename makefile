example:
	@echo "Hello World!"

.PHONY: example

gen-ent:
	go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/upsert,sql/modifier ./internal/adapter/ent/schema
add-gen-ent:
	go run -mod="mod" entgo.io/ent/cmd/ent new --target ./internal/adapter/ent/schema User
run:
	go run ./cmd/main.go