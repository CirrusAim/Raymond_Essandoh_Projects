//go:build tools

package quiz

import (
	_ "github.com/relab/gorums/cmd/protoc-gen-gorums"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
