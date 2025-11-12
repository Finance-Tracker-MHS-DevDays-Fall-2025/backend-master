
.PHONY: run
run:
	go run cmd/main.go

.PHONY: generate-http
generate-http:

.PHONY: protogen
protogen:
	protoc \
		--go_out=./internal/data/generated/ \
		--go-grpc_out=./internal/data/generated/
		./proto/service.proto
