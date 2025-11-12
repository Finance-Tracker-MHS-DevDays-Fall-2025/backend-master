
.PHONY: run
run:
	go run cmd/main.go

.PHONY: generate-http
generate-http:

.PHONY: protogen
protogen:
	buf generate ./internal/api/proto/
