# partly from https://medium.com/@olebedev/live-code-reloading-for-golang-web-projects-in-19-lines-8b2e8777b1ea

OS=linux
ARCH=386
BIN=bin/cavs-admin-$(OS)-$(ARCH)
NTVBIN=bin/cavs-admin
PID=/tmp/cavs-admin.pid

IMGREPO=tma1/cavs
VERSION=$(shell cat VERSION)

default: native
all: native linux

linux: $(BIN)
native: $(NTVBIN)

GO_FILES = $(wildcard ./**/*.go)

$(NTVBIN):
	go build -o $(NTVBIN) cmd/cavs-admin/main.go

$(BIN): $(wildcard ./**/*.go)
	GOOS=$(OS) GOARCH=$(ARCH) go build -o $(BIN) ./cmd/$(NTVBIN)

docker: $(BIN)
	docker build -t cavs-admin:$(VERSION) ./
	docker tag cavs-admin:$(VERSION) $(IMGREPO):$(VERSION)
	@echo 'Docker image built and tagged with $(VERSION)'
	@echo 'Next, run "docker push $(IMGREPO):$(VERSION)"'

install: $(NTVBIN) 
	cp $(NTVBIN) $(GOPATH)/bin/$(NTVBIN)

clean:
	if [[ -e $(NTVBIN) ]]; then rm $(NTVBIN); fi
	if [[ -e $(BIN) ]]; then rm $(BIN); fi

serve: restart
	@fswatch -o . -e ".*" -i "\\.go$$" -i "go.mod" | xargs -n1 -I{}  make restart || make kill

kill:
	@kill `cat $(PID)` || true

before:
	@echo "before noop"

test: $(GO_FILES)
	@echo $? $@

restart: kill clean before $(NTVBIN)
	@$(NTVBIN) & echo $$! > $(PID)

.PHONY: test build serve restart kill before docker install clean # let's go to reserve rules names

