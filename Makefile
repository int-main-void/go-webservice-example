#
#
#
#

all: image
IMAGE_DIR=image
PROJECT_NAME=example-project
VERSION=dev-latest

MAIN_PKG=greatsagemonkey.com/example-app
BIN_NAME=./bin/example-app
IMAGE_FILESET=bin pkg conf tools test

export GOPATH=$(PWD)

src/vendor: src/glide.yaml
	#go get github.com/gorilla/mux
	cd src/ && glide up

clean:
	rm -rf $(IMAGE_DIR)
	rm -rf bin/
	rm -rf pkg/

verify:
	go vet

.PHONY: test
test:
	go test -cover greatsagemonkey.com/...

bin,pkg: src/vendor test
	GOOS=linux GOARCH=386 go install $(MAIN_PKG)

image: bin,pkg
	mkdir -p $(IMAGE_DIR)/root
	cp conf/Dockerfile $(IMAGE_DIR)
	cp -a $(IMAGE_FILESET) $(IMAGE_DIR)/root
	docker build -t $(PROJECT_NAME):$(VERSION) $(IMAGE_DIR)

run: image
	docker run -it -v /opt/app/example-webservice/ssl:/opt/app/ssl -v /opt/app/example-webservice/log:/opt/app/log -e CONFIG_FILENAME="conf/example-webservice.conf" -e RUNTIME_STAGE="dev" $(PROJECT_NAME):$(VERSION)

bin-local: src/vendor test
	go install $(MAIN_PKG)

run-local: bin-local
	CONFIG_FILENAME="conf/example-webservice.conf" RUNTIME_STAGE="dev" $(BIN_NAME)
