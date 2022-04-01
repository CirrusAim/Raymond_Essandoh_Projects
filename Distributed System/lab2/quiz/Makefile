binaries := bin
participant_bin := $(binaries)/participant
quizmaster_bin := $(binaries)/quizmaster
gorum_include := $(shell go list -m -f {{.Dir}} github.com/relab/gorums)
grpc_proto := proto/base.proto proto/participant.proto
gorums_proto := proto/quizmaster.proto

all: pre quizmaster participant

.PHONY: pre
pre:
	@mkdir -p $(binaries)

quizmaster: protos
	@echo "+ compiling quizmaster server"
	@go build $(BUILD_FLAGS) -o $(quizmaster_bin) cmd/quizmaster/main.go

participant: protos
	@echo "+ compiling participant client"
	@go build $(BUILD_FLAGS) -o $(participant_bin) cmd/participant/main.go

.PHONY: tools
tools:
	@echo "+ installing tools"
	@go mod download
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -I % go install %

protos: tools grpcproto gorumsproto

grpcproto: $(grpc_proto)
	@echo "+ compiling grpc proto files"
	@protoc --go_out=paths=source_relative:. \
		--go-grpc_out=paths=source_relative:. $^

gorumsproto: $(gorums_proto)
	@echo "+ compiling gorums proto files"
	@protoc -I=$(gorum_include):. \
		--go_out=paths=source_relative:. \
		--gorums_out=paths=source_relative:. $^

clean:
	rm -vf $(participant_bin) $(quizmaster_bin)
