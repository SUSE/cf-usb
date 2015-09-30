# This makefile installs all dependencies (with the notable exception of the
# etcd binary) that are required for a CI build.

.PHONY: deps

deps:
	go get -u github.com/tools/godep
	go get -u github.com/jstemmer/go-junit-report
	cd ${GOPATH}/src/github.com/hpcloud/gocfbroker && godep restore
	cd ${GOPATH}/src/github.com/hpcloud/gocfbroker/boltdb && godep restore
	cd ${GOPATH}/src/github.com/hpcloud/gocfbroker/etcddb && godep restore
