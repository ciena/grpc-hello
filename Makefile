# Copyright 2021 Ciena Corporation.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
VERSION=$(shell head -1 ./VERSION)
DOCKER_REPOSITORY ?=
DOCKER_TAG ?= $(VERSION)
CLIENT_IMG ?= $(DOCKER_REPOSITORY)grpc-hello-client:$(DOCKER_TAG)
SERVER_IMG ?= $(DOCKER_REPOSITORY)grpc-hello-server:$(DOCKER_TAG)

BUILD_DATE=$(shell date -u "+%Y-%m-%dT%H:%M:%S%Z")
VCS_REF=$(shell git rev-parse HEAD)
ifeq ($(shell git ls-files --others --modified --deleted --exclude-standard | wc -l | tr -d ' '),0)
VCS_DIRTY=false
else
VCS_DIRTY=true
endif
ifeq ($(shell uname -s | tr '[:upper:]' '[:lower:]'),darwin)
VCS_COMMIT_DATE=$(shell date -j -u -f "%Y-%m-%d %H:%M:%S %z" "$(shell git show -s --format=%ci HEAD)" "+%Y-%m-%dT%H:%M:%S%Z")
else
VCS_COMMIT_DATE=$(shell date -u -d "$(shell git show -s --format=%ci HEAD)" "+%Y-%m-%dT%H:%M:%S%Z")
endif
# Remove any auth information from URL
GIT_TRACKING=$(shell git branch -vv | cut '-d ' -f 4 | sed -e 's/[][]//g' | cut -d/ -f1)
ifeq ($(GIT_TRACKING),)
GIT_TRACKING=origin
endif
VCS_URL=$(shell git remote get-url $(GIT_TRACKING) | sed -e 's/\/\/[-_:@a-zA-Z0-9]*[:@]/\/\//g')

DOCKER_BUILD_ARGS=\
--build-arg org_label_schema_version="$(VERSION)" \
--build-arg org_label_schema_vcs_url="$(VCS_URL)" \
--build-arg org_label_schema_vcs_ref="$(VCS_REF)" \
--build-arg org_label_schema_vcs_commit_date="$(VCS_COMMIT_DATE)" \
--build-arg org_label_schema_vcs_dirty="$(VCS_DIRTY)" \
--build-arg org_label_schema_build_date="$(BUILD_DATE)"

BUILD_LD_FLAGS=\
-X github.com/ciena/grpc-hello/pkg/version.version="$(VERSION)" \
-X github.com/ciena/grpc-hello/pkg/version.vcsURL="$(VCS_URL)" \
-X github.com/ciena/grpc-hello/pkg/version.vcsRef="$(VCS_REF)" \
-X github.com/ciena/grpc-hello/pkg/version.vcsCommitDate="$(VCS_COMMIT_DATE)" \
-X github.com/ciena/grpc-hello/pkg/version.vcsDirty="$(VCS_DIRTY)" \
-X github.com/ciena/grpc-hello/pkg/version.goVersion="$(shell go version 2>/dev/null | cut -d ' ' -f 3)" \
-X github.com/ciena/grpc-hello/pkg/version.os="$(shell go env GOHOSTOS)" \
-X github.com/ciena/grpc-hello/pkg/version.arch="$(shell go env GOHOSTARCH)" \
-X github.com/ciena/grpc-hello/pkg/version.buildDate="$(BUILD_DATE)"

.DEFAULT_GOAL:=help
.PHONY: help
help:  ## Display this help
	@echo "Usage: make \033[36m<target>\033[0m"
	@awk 'BEGIN {FS = ":.*## *"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "\033[36m  %s\033[0m,%s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST) | column -s , -c 2 -tx

pkg/apis/hello/hello.pb.go pkg/apis/hello/hello_grpc.pb.go: api/hello.proto
	protoc -I ./api --go_out=. --go-grpc_out=. $<

protos: pkg/apis/hello/hello_grpc.pb.go pkg/apis/hello/hello_grpc.pb.go

build: protos client server ## Build the local binaries for the client and server

client:  protos ## build local client binary
	go build -ldflags "$(BUILD_LD_FLAGS)" ./cmd/client

server: protos ## build local server binary
	go build -ldflags "$(BUILD_LD_FLAGS)" ./cmd/server

docker-build-client: ## Build docker image for client
	docker build $(DOCKER_BUILD_OPTIONS) . -t $(CLIENT_IMG) -f build/Dockerfile.client $(DOCKER_BUILD_ARGS)

docker-push-client: ## Push client docker image to repository
	docker push $(DOCKER_PUSH_OPTIONS) $(CLIENT_IMG)

docker-build-server: ## Build docker image for server
	docker build $(DOCKER_BUILD_OPTIONS) . -t $(SERVER_IMG) -f build/Dockerfile.server $(DOCKER_BUILD_ARGS)

docker-push-server: ## Push server docker image to repository
	docker push $(DOCKER_PUSH_OPTIONS) $(SERVER_IMG)

docker-build: docker-build-client docker-build-server ## Build docker images for client and server

docker-push: docker-push-client docker-push-server ## Push docker images to repository

clean:
	rm -rf client server
