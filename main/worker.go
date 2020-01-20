package main

import (
	"cs675-spring20-labs/lab1/serverless"
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"plugin"
	"sync"
)

// Worker holds the state for a server waiting for:
// 1) RegisterService,
// 2) InvokeService,
// 3) Shutdown RPCs.
type Worker struct {
	sync.Mutex

	address    string
	masterAddr string
	nThreads   int
	nTasks     int
	concurrent int
	l          net.Listener

	shutdown chan struct{}
}

// Service struct maintains the state of a plugin Service.
type Service struct {
	pluginDir string
	name      string
	interf    serverless.Interface
}

// serviceMap is a global map that keeps track of all registered plugin Services.
// You should insert the newly registered service into this map.
var serviceMap = make(map[string]*Service)

// newService initializes a new plugin Service.
func newService(serviceName string) *Service {
	return &Service{
		pluginDir: serverless.PluginDir,
		name:      serviceName,
		interf:    nil,
	}
}

// RegisterService is caled by the driver to plugin a new service that has already been
// compiled into a .so static object library.
func (wk *Worker) RegisterService(args *serverless.ServiceRegisterArgs, _ *struct{}) error {
	plug, err := plugin.Open("../plugins/" + args.ServiceName + ".so")
	if err != nil {
		fmt.Printf("Failed to open plugin %s: %v\n", args.ServiceName, err)
		return err
	}
	// TODO: implement me
	// Hint 1: You may want to use `plug.Lookup` to locate the service symbol,
	// and expose the interested service API associated with serverless.Interface.
	// Hint 2: Call newService to initialize a service struct, and insert the service key-value pair
	// to the global serviceMap.
	// TODO TODO TODO
	//

	serverless.Debug("Successfully registered new service %s\n", args.ServiceName)
	return nil
}

// InvokeService is called by the driver (schedule) when a new task
// is being scheduled on this worker.
func (wk *Worker) InvokeService(args serverless.RPCArgs, _ *struct{}) error {
	// TODO: implement me
	// Hint: You should locate the interested service registered from serviceMap.
	// and call `service.interf.DoService` to make the call to the plugin service.
	// TODO TODO TODO
	//
	return nil
}

// Shutdown is called by the driver when all work has been completed.
// No response needed.
func (wk *Worker) Shutdown(_ *struct{}, _ *struct{}) error {
	serverless.Debug("Worker shutdown %s\n", wk.address)
	close(wk.shutdown)
	wk.l.Close()
	return nil
}

// Tell the driver I exist and ready to work:
// register is the internal function that calls the RPC method of Driver.Register
// at the remote driver to register the worker itself.
func (wk *Worker) register(driver string) {
	args := new(serverless.WorkerRegisterArgs)
	args.WorkerAddr = wk.address

	ok := serverless.Call(driver, "Driver.Register", args, new(struct{}))
	if ok == true {
		fmt.Printf("Successfully registered worker %s\n", wk.address)
	} else {
		fmt.Printf("Failed to register worker %s\n", wk.address)
	}
}

// startRPCServer sets up a connection with the driver, registers its address,
// and waits for any of the following two events:
// 1) plugin Services to be registered,
// 2) tasks to be scheduled.
func (wk *Worker) startRPCServer() {
	// TODO: implement me
	// Hint: Refer to how the driver's startRPCServer is implemented.
	// TODO TODO TODO
	//

	//
	// Once shutdown is closed, should the following statement be
	// called, meaning the worker RPC server is existing.
	serverless.Debug("Worker: %v RPC server exist\n", wk.address)
}

// The main entrance
func main() {
	wk := new(Worker)
	wk.address = os.Args[1]    // the 1st cmd-line argument: worker hostname and ip addr
	wk.masterAddr = os.Args[2] // the 2nd cmd-line argument: driver hostname and ip addr
	wk.shutdown = make(chan struct{})

	wk.startRPCServer()
}
