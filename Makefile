
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
	docker build . -f build/Dockerfile -t master --load

.PHONY: upload-img
upload-img:
	docker tag master:latest cr.yandex/crpkimlhn85fg9vjfj7l/master:latest
	docker image push cr.yandex/crpkimlhn85fg9vjfj7l/master:latest
