package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// To compile the plugin: run:
// go build --buildmode=plugin -o helloworkd_service.so helloworld_service.go

// Define helloWorldService
type helloWorldService string

// HelloWorldArgs defines this plugin's argument format
type HelloWorldArgs struct {
	TaskNum int // just a single field indicating the task number
}

// DoService is called remotely by the scheduler (essentially the
// driver) when a new task is being scheduled and invoked on this
// worker.
func (s helloWorldService) DoService(raw []byte) error {
	var args HelloWorldArgs
	// TODO: implement me
	// Hint 1: You may want to use `bytes.NewBuffer` and
	// `gob.NewDecoder` etc. to deserialize (unmarshal) the raw
	// []byte array of the argument.
	// Hint 2: When successfully deserializing the argument, you
	// should call `fmt.Printf` to print out "Hello world from
	// helloWorldService plugin" together with the task number you
	// parsed from the raw []byte array.
	// TODO TODO TODO
	//

	return nil
}

// Expose the helloWorldService interface to the worker
var Interface helloWorldService
