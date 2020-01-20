package serverless

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"sync"
)

// schedule starts and waits for all tasks to finish.
// In Lab 1, this is just to invoke the newly registered plugin
// service at the worker side.
func (drv *Driver) schedule(
	registerChan chan string,
	serviceName string,
) {
	// HelloWorldArgs should follow the argument format defined by
	// your helloworld_service plugin.
	type HelloWorldArgs struct {
		TaskNum int
	}

	// For Lab 1, the number of tasks is hard-coded. This won't be
	// the case for Lab 2.
	nTasks := 10

	// readyChan is a bounded buffer that is used to notify the
	// scheduler of workers that are *TRULY* ready for executing the
	// service tasks.
	readyChan := make(chan string, nTasks)

	// invokeService is a goroutine that is used to call the RPC
	// method of Worker.InvokeService at the worker side.
	invokeService := func(worker string, args *HelloWorldArgs) {
		// TODO: implement me
		// Hint: You may find the `gob.NewEncoder` useful for
		// serializing (marshalling) the service parameters (in this
		// case: HelloWorldArgs struct).
		// TODO TODO TODO
		//

		// Notify the scheduler that this worker is back to ready state.
		readyChan <- worker
		if !success {
			fmt.Printf("Schedule: task failed to execute by %v: %v\n", worker, args.TaskNum)
			os.Exit(1)
		}
	}

	// TODO: implement me
	// All tasks have to be scheduled on workers, and only once all of them
	// have been completed successfully should the function return.
	// Hint 1: Use for loop to loop over all the tasks, and use select
	// inside of for loop between registerChan and readyChan.
	// Hint 2: You may want to use goroutine to invoke function invokeService
	// like the following:
	// go invokeService(worker, args)
	// TODO TODO TODO
	//

	// Work done: finish the task scheduling
	Debug("Driver: Task scheduling done\n")
}
