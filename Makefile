
.PHONY: run
run:
	go run cmd/main.go

.PHONY: protogen
protogen:
	protoc \
		--go_out=./internal/data/generated/ \
		--go-grpc_out=./internal/data/generated/
		./proto/service.proto
