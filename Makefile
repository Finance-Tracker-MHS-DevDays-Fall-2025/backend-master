
.PHONY: run
run:
	go run cmd/main.go

.PHONY: protogen
protogen:
	buf generate ./internal/api/proto/
