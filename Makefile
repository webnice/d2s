## Simple projects tooling for every day
## Версия: 2023.12.13

DIR                   = $(strip $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))
GOPATH               := $(GOPATH)
DATE                 := $(shell date -u +%Y%m%d.%H%M%S.%Z)
GOGENERATE            = $(shell if [ -f .gogenerate ]; then cat .gogenerate; fi)
TESTPACKETS           = $(shell if [ -f .testpackages ]; then cat .testpackages; fi)
BENCHPACKETS          = $(shell if [ -f .benchpackages ]; then cat .benchpackages; fi)
GO111MODULE          ?= $(GO111MODULE:on)
PACKAGES_LOCK_VER     = $(shell if [ -f .packages_lock_version ]; then cat .packages_lock_version; fi)
COVERALLS_REPO_TOKEN  = $(shell if [ -f .coveralls ]; then cat .coveralls; fi)

default: help

link:
.PHONY: link

## Кодогенерация.
generate:
.PHONY: generate

## Загрузка зависимостей.
dep-init:
	@rm -rf ${DIR}/vendor 2>/dev/null; true
.PHONY: dep-init
dep: dep-init
	@go mod download
	@go mod tidy
	@go mod vendor
.PHONY: dep

update: upd
upd:
	@go clean -cache -modcache
	@go get -u ./...
	@for item in $(PACKAGES_LOCK_VER); do \
		go get -u "$${item}"; \
		true; \
	done
.PHONY: upd
.PHONY: update

## Кодогенерация (run only during development).
## All generating files are included in a .gogenerate file.
gen:
	@for PKGNAME in $(GOGENERATE); do GOPATH="$(DIR)" go generate $${PKGNAME}; done
.PHONY: gen

build:
	@go build -o bin/db2struct github.com/webnice/d2s/v2/db2struct
.PHONY: build

## Тестирование одного или нескольких пакетов, а также приложений с отчётом о проценте покрытия тестами.
## Все тестовые файлы включены в файл .testpackages.
test: link
	@echo "mode: count" > $(DIR)/coverage.log
	@for PACKET in $(TESTPACKETS); do \
		touch $(DIR)/coverage-tmp.log; \
		unset GOPATH; go test -v -covermode=count -coverprofile=$(DIR)/coverage-tmp.log $$PACKET || exit 1; \
		if [ "$$?" -ne "0" ]; then exit $$?; fi; \
		tail -n +2 $(DIR)/coverage-tmp.log | sort -r | awk '{if($$1 != last) {print $$0;last=$$1}}' >> $(DIR)/coverage.log; \
		rm -f $(DIR)/coverage-tmp.log; true; \
	done
	@if [ "$(COVERALLS_REPO_TOKEN)" != "" ]; then\
		echo "Загрузка в coveralls.io процента покрытия кода тестами."; \
		goveralls -repotoken $(COVERALLS_REPO_TOKEN) >/dev/null 2>&1; true; \
	fi
.PHONY: test

## Отображение в браузере покрытия тестируемого кода в HTML-отчёте (запускается только во время разработки).
cover: test
	unset GOPATH; go tool cover -html=$(DIR)/coverage.log
.PHONY: cover

## Тестирование производительности.
## Все тестовые файлы включены в файл .benchpackages.
bench:
	@for PACKET in $(BENCHPACKETS); do GOPATH=${GOPATH} go test -race -bench=. -benchmem $$PACKET; done
.PHONY: bench

## Очистка от временных файлов.
clean:
	@if [ -d ${DIR}/pkg ]; then\
		chown -R `whoami` ${DIR}/pkg/; true; \
		chmod -R 0777 ${DIR}/pkg/; true; \
	fi
	@rm -rf ${DIR}/src; true
	@rm -rf ${DIR}/bin; true
	@rm -rf ${DIR}/pkg; true
	@rm -rf ${DIR}/rpmbuild; true
	@rm -rf ${DIR}/*.log; true
	@rm -rf ${DIR}/*.lock; true
.PHONY: clean

## Помощь по командам.
help:
	@echo "Usage: make [target]"
	@echo "  target is:"
	@echo "    gen                  - Кодогенерация с использованием go generate."
	@echo "    test                 - Запуск тестов проекта."
	@echo "    cover                - Запуск тестов проекта с отображением процента покрытия кода тестами."
	@echo "    bench                - Запуск тестов производительности."
	@echo "    clean                - Очистка папки проекта от временных файлов."
.PHONY: help
