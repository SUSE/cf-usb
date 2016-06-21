#!/usr/bin/env make

GIT_ROOT:=$(shell git rev-parse --show-toplevel)

.PHONY: all clean format lint vet build test tools dist

all: clean format lint vet bindata build test

clean:
	${GIT_ROOT}/make/clean

format:
	${GIT_ROOT}/make/format

lint:
	${GIT_ROOT}/make/lint

vet:
	${GIT_ROOT}/make/vet

build: 
	${GIT_ROOT}/make/build

test:
	${GIT_ROOT}/make/test

tools:
	${GIT_ROOT}/make/tools


genswagger:
	${GIT_ROOT}/make/genswagger
	
dist:  
	${GIT_ROOT}/make/dist