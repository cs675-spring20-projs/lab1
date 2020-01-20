package main

import (
	"cs675-spring20-labs/lab1/serverless"
	"os"
)

// The client can be run in the following way:
// go run client.go localhost:1234 helloworld_service
// where localhost:1234 is the driver (client)'s hostname and ip address,
// and helloworld_service is the name of the service plugin you want to execute on remote worker
func main() {
	drv := serverless.NewDriver(os.Args[1]) // the 1st cmd-line argument: driver hostname and ip addr
	serviceName := os.Args[2]               // the 2nd cmd-line argument: plugin Service name

	go drv.Run(serviceName)

	drv.Wait()
}
