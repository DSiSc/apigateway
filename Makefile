# Copyright(c) 2018 DSiSc Group All Rights Reserved.
# 
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.PHONY: default help all build test devenv gotools clean

VERSION=$(shell grep "const Version" version/version.go | sed -E 's/.*"(.+)"$$/\1/')
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')

default: all

help:
	@echo 'Management commands for DSiSc/apigateway:'
	@echo
	@echo 'Usage:'
	@echo '    make build           Compile the project.'
	@echo '    make vet-test        Run vet tests on a compiled project.'
	@echo '    make test            Run tests on a compiled project.'
	@echo '    make devenv          Prepare devenv for test or build.'
	@echo '    make fetch-deps      Run govendor fetch for deps.'
	@echo '    make gotools         Prepare go tools depended.'
	@echo '    make clean           Clean the directory tree.'
	@echo

all: test build

build:
	echo "building apigateway ${VERSION}"
	echo "GOPATH=${GOPATH}"
	go build -ldflags "-X github.com/DSiSc/apigateway/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X github.com/DSiSc/apigateway/version.BuildDate=${BUILD_DATE}" ./...

vet-test:
	go vet `go list ./...`

test: vet-test 
	cd unit-test && ./run.sh

## tools & deps
devenv: gotools fetch-deps

fetch-deps: gotools-clean
	govendor init && govendor fetch +m 

gotools:
	cd gotools && make install

.PHONY: clean
clean: gotools-clean

gotools-clean:
	cd gotools && make clean
