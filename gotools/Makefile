# Copyright(c) 2018 DSiSc Group. All Rights Reserved.
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

BINDIR ?= /usr/local/bin
OBJDIR ?= build

TMP_GOPATH=$(OBJDIR)/gopath
GOBIN=$(abspath $(TMP_GOPATH)/bin)

GOTOOLS = gotype goimports godoc gorename gomvpkg golint \
             govendor gocov misspell

GOTOOLS_BIN = $(patsubst %,$(GOBIN)/%, $(GOTOOLS))

# go tool->path mapping
## STD_TOOLS
go.fqp.gotype    := golang.org/x/tools/cmd/gotype
go.fqp.goimports := golang.org/x/tools/cmd/goimports
go.fqp.godoc     := golang.org/x/tools/cmd/godoc
go.fqp.gorename  := golang.org/x/tools/cmd/gorename
go.fqp.gomvpkg   := golang.org/x/tools/cmd/gomvpkg
go.fqp.golint    := github.com/golang/lint/golint

## EXTRA_TOOLS
go.fqp.govendor  := github.com/kardianos/govendor
go.fqp.gocov     := github.com/axw/gocov/...
go.fqp.misspell  := github.com/client9/misspell/cmd/misspell
all: $(GOTOOLS_BIN)

install: $(GOTOOLS_BIN)
	@mkdir -p $(BINDIR)
	cp $^ $(BINDIR)

# Default rule for gotools uses the name->path map for a generic 'go get' style build
gotool.%:
	$(eval TOOL = ${subst gotool.,,${@}})
	@echo "Building ${go.fqp.${TOOL}} -> $(TOOL)"
	@GOPATH=$(abspath $(TMP_GOPATH)) go get -v ${go.fqp.${TOOL}}

$(GOBIN)/%:
	$(eval TOOL = ${subst $(GOBIN)/,,${@}})
	@$(MAKE) gotool.$(TOOL)

.PHONY: clean
clean:
	@rm -rf $(OBJDIR)
