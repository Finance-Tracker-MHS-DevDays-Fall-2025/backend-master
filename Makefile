
.PHONY: env
env:
	cp .env.example .env

.PHONY: sync-submodules
sync-submodules:
	git submodule update --init --recursive --remote --merge

.PHONY: protogen
protogen:
	buf generate ./internal/api/proto/

.PHONY: run
run:
	go run cmd/main.go

.PHONY: build-img
build-img:
	docker build . -f build/Dockerfile -t master-service

.PHONY: upload-img
upload-img:
	docker tag master-service:latest cr.yandex/crpkimlhn85fg9vjfj7l/master-service:latest
	docker image push cr.yandex/crpkimlhn85fg9vjfj7l/master-service:latest
