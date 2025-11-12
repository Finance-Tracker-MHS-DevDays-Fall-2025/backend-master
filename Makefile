
.PHONY: sync-submodules
sync-submodules:
	git submodule update --init --recursive --remote --merge

.PHONY: protogen
protogen:
	buf generate ./internal/api/proto/

.PHONY: run
run:
	go run cmd/main.go
