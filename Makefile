NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
DEPS=$(shell go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
GOPKGS=$(shell go list -f '{{.Dir}}' ./...)
OSES=linux darwin

include version.mk

COMMIT_HASH=$(shell git log --pretty=format:'%h' -n 1)
APP_VERSION=$(VERSION)-$(COMMIT_HASH)
DIST_FIND_BUILDS=find * -type d -not -path "forpatches" -exec

.PHONY: all dist format lint vet build test tools bench clean generate cleangeneratedfiles
.SILENT: all dist format lint vet build test tools bench clean generate cleangeneratedfiles

all: clean format lint vet build test dist

format:
	@echo "$(OK_COLOR)==> Checking format$(ERROR_COLOR)"
	@echo $(PKGSDIRS) | tr ' ' '\n' | xargs -I '{p}' -n1 goimports -e -l {p} | sed "s/^/Failed: /"
	@echo "$(NO_COLOR)\c"

lint:
	@echo "$(OK_COLOR)==> Linting$(ERROR_COLOR)"
	@echo $(PKGSDIRS) | tr ' ' '\n' | xargs -I '{p}' -n1 golint {p} | grep -v "mock_.*\.go" | sed "s/^/Failed: /"
	@echo "$(NO_COLOR)\c"

vet:
	@echo "$(OK_COLOR)==> Vetting$(ERROR_COLOR)"
# Blame https://code.google.com/p/go/issues/detail?id=6820 for the -composites=false
	@echo $(GOPKGS) | tr ' ' '\n' | xargs -I '{p}' -n1 go tool vet -composites=false {p} | sed "s/^/Failed: /"
	@echo "$(NO_COLOR)\c"

driversbindata:
	@echo "$(OK_COLOR)==> Embedding JSON schemas into drivers$(NO_COLOR)"
	find driver/ -maxdepth 1 -type d \( ! -name driver \) -exec \
	bash -c "(cd '{}' && go-bindata -pkg="driverdata" -o driverdata/schemas.go schemas/ )" \;

build: generate build-usb build-drivers build-driver-generator
	#       calls all the other necessary builds

build-drivers : build-dummy-async-driver build-dummy-driver build-mongo-driver build-mysql-driver

build-drivers : build-mssql-driver build-postgres-driver build-rabbitmq-driver build-redis-driver

build-usb:
	$(call buildme,./cmd/usb,usb)


build-driver-generator:
	$(call buildme,./driver-generator/cmd/driver-generator,driver-generator)

build-dummy-async-driver:
	$(call buildme,./cmd/driver/dummy-async,drivers)

build-dummy-driver:
	$(call buildme,./cmd/driver/dummy,drivers)

build-mongo-driver:
	$(call buildme,./cmd/driver/mongo,drivers)

build-mssql-driver:
	$(call buildme,./cmd/driver/mssql,drivers)

build-postgres-driver:
	$(call buildme,./cmd/driver/postgres,drivers)

build-rabbitmq-driver:
	$(call buildme,./cmd/driver/rabbitmq,drivers)

build-redis-driver:
	$(call buildme,./cmd/driver/redis,drivers)

build-mysql-driver:
	$(call buildme,./cmd/driver/mysql,drivers)


buildme =  @(DIRNAME=$$(basename $(1));\
	echo "$(OK_COLOR)==> Building$(NO_COLOR) $$DIRNAME"; \
	for OS in $(OSES); do \
		DIRNAME=$$(basename $(1)); \
		 env GOOS=$$OS GOARCH=amd64 go build \
		 -ldflags="-X main.version=$(APP_VERSION)" \
		 -o build/$$OS-amd64/$(2)/$$DIRNAME $(1); \
	done)

patch: cleangeneratedfiles
	@echo "$(OK_COLOR)==> Generating Update Patches$(NO_COLOR)"
	export CGOENABLED=1 && \
	export GOPATH=$(shell godep path):$(shell echo $$GOPATH) &&\
	gox -verbose \
	-ldflags="-X main.version $(APP_VERSION) -X main.updateserver $(PATCH_SERVER) -X main.branch $(PATCH_CHANNEL)" \
	-os="windows linux darwin " \
	-arch="amd64" \
	-output="build/forpatches/{{.OS}}-{{.Arch}}" ./...
	mv build/forpatches/windows-amd64.exe build/forpatches/windows-amd64

	go-selfupdate "build/forpatches" $(APP_VERSION)
	ln -s public cf-mgmt_$(MAJOR_MINOR)_$(PATCH_CHANNEL)

cleangeneratedfiles:
	rm -rf build/forpatches
	rm -rf public

test:
	@echo "$(OK_COLOR)==> Testing$(NO_COLOR)"
	gocov test $(shell go list ./... | grep -v /vendor/) -v -timeout 60m | gocov-xml > coverage.xml
	@echo "$(NO_COLOR)\c"

tools:
	@echo "$(OK_COLOR)==> Installing tools$(NO_COLOR)"
	#Great tools to have, and used in the build file
	go get -u golang.org/x/tools/cmd/goimports
	go get -u golang.org/x/tools/cmd/cover
	go get -u github.com/golang/lint/golint
	#Tools for the project
	go get -u github.com/codegangsta/cli
	go get -u github.com/tools/godep
	go get -u github.com/mitchellh/gox
	go get -u github.com/vektra/mockery/cmd/mockery
	go get -u github.com/axw/gocov/gocov
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	go get -u github.com/jteeuwen/go-bindata/...
	
	# gox -build-toolchain
	#dependencies for project
	go get gopkg.in/yaml.v2
	go get github.com/stretchr/testify/assert
	go get github.com/stretchr/testify/mock
	#Tools for code coverage reporting
	go get github.com/axw/gocov/...
	go get github.com/AlekSi/gocov-xml
	go get gopkg.in/matm/v1/gocov-html
	#Tools for integration tests
	go get github.com/nats-io/gnatsd
	#Fix for consul
	go get github.com/hashicorp/consul

clean: cleangeneratedfiles
	@echo "$(OK_COLOR)==> Cleaning$(NO_COLOR)"
	rm -rf build
	rm -rf $(GOPATH)/pkg/*
	rm -f $(GOPATH)/bin/usb

genswagger:
	@echo "$(OK_COLOR)==> Generationg management APIs using swagger$(NO_COLOR)"
	rm -rf lib/operations lib/genmodel
	swagger generate server -f swagger-spec/api.json -m genmodel -s "" -A usb-mgmt -t lib
	rm -rf lib/cmd
	go-bindata -pkg="data" -o lib/data/swagger.go swagger-spec/

dist: build 
ifeq ("$(PLATFORM)","") 
	@echo "$(OK_COLOR)==> Disting all$(NO_COLOR)"; \
	for OS in $(OSES); do \
		cd build/$$OS-amd64/ 1> /dev/null; tar czf ../../cf-usb-$(APP_VERSION)-$$OS-amd64.tgz ./; cd - 1> /dev/null; \
	done; 
else 
	@echo "$(OK_COLOR)==> Disting $(PLATFORM)$(NO_COLOR)"; \
	cd build/$(PLATFORM)-amd64/ 1> /dev/null; tar czf ../../cf-usb-$(APP_VERSION)-$(PLATFORM)-amd64.tgz ./; cd - 1> /dev/null; 
endif
