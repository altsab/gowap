PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)

export GO111MODULE=on

TEST_FLAGS := "" # "-mod=vendor"

default: help

getdeps:
	@echo "Installing golint" && go get -u golang.org/x/lint/golint

verifiers: lint

gen:
	@go generate

fmt:
	@echo "Running $@"
	@${GOPATH}/bin/golint .

lint:
	@echo "Running $@"
	@${GOPATH}/bin/golint -set_exit_status ./...

metalinter:
	@${GOPATH}/bin/gometalinter --install
	@${GOPATH}/bin/gometalinter --disable-all \
		-E vet \
		-E gofmt \
		-E misspell \
		-E ineffassign \
		-E goimports \
		-E deadcode --tests --vendor ./...

check: verifiers test

test:
	@echo "Running unit tests"
	@go test -v $(TEST_FLAGS) -tags kqueue ./...

bench:
	@echo "Running bench"
	@go test -bench=. -benchmem -benchtime=5s ./...

coverage:
	@echo "Running all coverage for surferua"
	@(env bash $(PWD)/go-coverage.sh)

pkg-add:
	@echo "Adding new package $(PKG)"
	@${GOPATH}/bin/govendor add $(PKG)

pkg-update:
	@echo "Updating new package $(PKG)"
	@${GOPATH}/bin/govendor update $(PKG)

pkg-remove:
	@echo "Remove new package $(PKG)"
	@${GOPATH}/bin/govendor remove $(PKG)

pkg-list:
	@$(GOPATH)/bin/govendor list

clean:
	@echo "Cleaning up all the generated files"
	@find . -name '*.test' | xargs rm -fv
	@rm coverage.txt
	@rm -rvf build
	@rm -rvf release

help:
	@echo "nothing to do!"
