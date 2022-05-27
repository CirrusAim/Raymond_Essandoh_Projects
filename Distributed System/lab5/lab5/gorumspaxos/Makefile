binaries := bin
paxosclient_bin := $(binaries)/paxosclient
paxosserver_bin := $(binaries)/paxosserver
ifeq ($(OS),Windows_NT)
    paxosclient_bin = $(binaries)/paxosclient.exe
    paxosserver_bin = $(binaries)/paxosserver.exe
endif
gorum_include := $(shell go list -m -f {{.Dir}} github.com/relab/gorums)
proto_src := proto/multipaxos.proto fdproto/fd.proto
proto_go := $(proto_src:%.proto=%.pb.go)

all: pre protos generateserver generateclient

.PHONY: pre
pre:
	@mkdir -p $(binaries)

.PHONY: protos
protos: $(proto_go)
	@echo "+ compiling gorums proto files"

%.pb.go %_gorums.pb.go : %.proto
	@protoc -I=$(gorum_include):. \
		--go_out=paths=source_relative:. \
		--gorums_out=paths=source_relative:. \
		$<

generateserver:
	@echo "+ compiling paxos server "
	@cp -f cmd/paxosserver/runserver.sh $(binaries)
	@go build $(BUILD_FLAGS) -o $(paxosserver_bin) cmd/paxosserver/main.go

generateclient:
	@echo "+ compiling paxos client "
	@cp -f cmd/paxosclient/runclient.sh $(binaries)
	@cp -f cmd/paxosclient/bench.sh $(binaries)
	@go build $(BUILD_FLAGS) -o $(paxosclient_bin) cmd/paxosclient/main.go

.PHONY: clean
clean:
	rm -rf $(binaries)/*