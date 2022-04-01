# Lab 2: Network Programming in Go

| Lab 2: | Network Programming in Go |
| ---------------------    | --------------------- |
| Subject:                 | DAT520 Distributed Systems |
| Deadline:                | **February 11, 2022 23:59** |
| Expected effort:         | 20-30 hours |
| Grading:                 | Pass/fail |
| Submission:              | Individually |

## Table of Contents

1. [Table of Contents](#table-of-contents)
2. [Introduction](#introduction)
3. [UDP Echo Server](#udp-echo-server)

## Table of Contents

1. [Introduction](https://github.com/dat520-2022/assignments/tree/main/lab2/README.md#introduction)
2. [UDP Echo Server](https://github.com/dat520-2022/assignments/tree/main/lab2/README.md#udp-echo-server)
3. [GRPC Server](grpc/README.md)
4. [Quizzing platform with gRPC and Gorums](quiz/quiz.md)
5. [Networking Questions](networking/network_questions.md)

## Introduction

The goal of this lab assignment is to get you started with network programming
in Go. The overall aim of the lab project is to implement a fault tolerant
distributed application.
Knowledge of network programming in Go is naturally a prerequisite for accomplishing this.

This lab assignment consist of four parts. In the first part you are expected to
implement a simple echo server that is able to respond to different commands
specified as text. In the second part, you will be implementing an in-memory
key-value storage using the `GRPC` framework as described [here](grpc/README.md).
You will need this latter skill, to implement the third part of this assignment
which consist of a quizzing platform with `GRPC` and [Gorums](https://github.com/relab/gorums) as
described [here](quiz/quiz.md).

The most important package in the Go standard library that you will use in
the first part of this assignment is the [`net`](http://golang.org/pkg/net) package.
It is recommended that you actively use the documentation available for this package during your
work on this lab assignment.

Finally, the last part of this lab consists of answering the [questions](networking/network_questions.md)
related to the knowledge of networking programming you obtained while implementing the previous parts.

## UDP Echo Server

In this task we will focus on the user datagram protocol (UDP), which provides
unreliable datagram service. You will find the documentation of the
[UDPConn](https://golang.org/pkg/net/#UDPConn) type useful.

In the provided code under `uecho`, we have implemented a simple
`SendCommand()` function that acts as a client, along with a bunch of tests.
You can run these test with `go test -v`, and as described in Lab 1, you can
use the `-run` flag to run only a specific test.

You can also compile your server code into a binary using `go build`. This
will produce a file named `uecho` in the same folder as the `.go` source files.
You can run this binary in two ways:

1. `./uecho -server &` will start the server in the background. Note: *This will
   not work until you have implemented the necessary server parts.*

2. `./uecho` will start the command line client, from which you may interact with
   the server by typing commands into the terminal window.

If you want to extend the capabilities of this runnable client and server,
you can edit the files `echo.go` and `echo_client.go`. But note that the
tests executed by the quickfeed will use original `SendCommand()` provided
in the original `echo_client.go` file. If you've done something fancy,
and want to show us that's fine, but it won't be considered by the quickfeed.

#### Echo server specification:


The `SendCommand()` takes the following arguments:

| Argument | Description	|
| -------------------- 	| ------------------------------------- |
| `udpAddr`		| UDP address of the server (`localhost:12110`) 		|
| `cmd`			| Command (as a text string) that the server should interpret and execute |
| `txt`			| Text string on which the server should perform the command provided in `cmd` |

The `SendCommand()` function produces a string composed of the following

```
cmd|:|txt
```

For example:

```
UPPER|:|i want to be upper case
```

From this, the server is expected to produce the following reply:

```
I WANT TO BE UPPER CASE
```

See below for more details about the specific behaviors of the server.

1. For each of the following commands, implement the corresponding functions, so that the returned value corresponds to the expected test outcome. Here you are expected to implement demultiplexer that demultiplexes the input (the command) so that different actions can be taken. A hint is to use the `switch` statement. You will probably also need the `strings.Split()` function.

    | Command	| Action |
    | -------------------- 	| ------------------------------------- |
    | UPPER		| Takes the provided input string `txt` and applies the translates it to upper case using `strings.ToUpper()`. |
    | LOWER		| Same as UPPER, but lower case instead. |
    | CAMEL		| Same as UPPER, but title or camel case instead. |
    | ROT13		| Takes the provided input string `txt` and applies the rot13 translation to it; see lab1 for an example. |
    | SWAP		| Takes the provided input string `txt` and inverts the case. For this command you will find the `strings.Map()` function useful, together with the `unicode.IsUpper()` and `unicode.ToLower()` and a few other similar functions. |

2. The server should reply `Unknown command` if it receives an unknown command
   or fails to interpret a request in any way.

3. Make sure that your server continues to function even if one client's
   connection or datagram packet caused an error.

#### Echo server implementation

You should implement the specification by extending the skeleton code found in `echo_server.go`:

```go
package main

import (
	"net"
	"strings"
)

// UDPServer implements the UDP server specification found at
// https://github.com/dat520-2022/assignments/tree/main/lab2/README.md#udp-server
type UDPServer struct {
	conn *net.UDPConn
	// TODO(student): Add fields if needed
}

// NewUDPServer returns a new UDPServer listening on addr. It should return an
// error if there was any problem resolving or listening on the provided addr.
func NewUDPServer(addr string) (*UDPServer, error) {
	// TODO(student): Implement
	return nil, nil
}

// ServeUDP starts the UDP server's read loop. The server should read from its
// listening socket and handle incoming client requests as according to the
// the specification.
func (u *UDPServer) ServeUDP() {
	// TODO(student): Implement
}

// socketIsClosed is a helper method to check if a listening socket has been
// closed.
func socketIsClosed(err error) bool {
	if strings.Contains(err.Error(), "use of closed network connection") {
		return true
	}
	return false
}
```