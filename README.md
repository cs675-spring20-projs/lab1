# CS 675 Lab 1: Serverless and RPC

## Important dates

**Due** Friday, 02/14, midnight.

This is also an **individual** lab.

## Change log

* **02/02/20:** Updated how to setup `$GOPATH`.
* **02/02/20:** Fixed typos. 

## Introduction

In Lab 1, you will build a client and server framework as a way to
to learn the Go programming language and as a way to learn about RPCs
in distributed systems. 

## Software

You'll implement this Lab (and all the Labs) in  <a
href="http://www.golang.org/">Go</a>. The Go website contains lots
of tutorial information which you may want to look at.

The second half of the [Environment Setup instructions](https://tddg.github.io/cs675-spring20/env_setup.html)
lists a fairly good amount of Go resources (Go editors, coding style, useful tools, etc.). Definitely take a 
look and give them a try.

In this Lab, we supply you with parts of a flexible serverless Go
RPC framework implementation. We call it a <a
href="https://aws.amazon.com/lambda/">*serverless*</a> framework, not
because *serverless* is hot, but because we will use Go's plugin
feature to implement a dynamically pluggable RPC framework that at
its core realizes a serverless model -- the user (i.e., the developer
of the Lambda function, or in our case, the Go plugin service
modules) don't need to worry about servers but can simply focus on
the development of the core business logic -- her own plugin Lambda
function.

### Getting familiar with the source

The serverless package (located at `serverless`) provides
a simple serverless library with a partially implemented RPC plugin
framework. Applications would normally call `Run()` located in
`serverless/driver.go` to register a Lambda function (i.e., a 
Go plugin service library) and start a configurable number of tasks
(executing the registered Lambda function).

The flow of the RPC client-server implementation is as follows:
		
1. The user (i.e., you) provides one or multiple Lambda functions in
the form of a Go plugin library (see the skeleton code
`helloworld_service.go` provided under `plugins/`).
		
2. A driver is created with this knowledge when running
`client` (see `main/client.go`). It spins up an RPC
server (see `serverless/driver.go`), and waits for workers to
register (using the RPC call `Register()` defined in
`serverless/driver.go`). 
		
3. One or multiple worker processes are created when running
`worker` (see `main/worker.go`). Each worker spins up
an RPC server, and registers itself at the driver, and waits for the
driver to register Lambda functions and schedule tasks.
		
4. The driver registers the Lambda function specified as
`serviceName` (the second command-line argument provided when running
`driver.go`) by calling `prepareService()` (see
`serverless/driver.go`), which issues the
`RegisterService()` RPC call to a worker.  The driver then
schedules the tasks by calling `schedule()` (see
`serverless/schedule.go`). In `schedule()`, the driver
issues the `InvokeService()` RPC call to execute an
already-registered plugin Lambda function on a worker.
		
4. The driver sends a `Shutdown()` RPC to each of its
workers, and then shuts down its own RPC server.

You should look through the files in the whole framework
implementation, as reading them would be useful to understand how the
other methods fit into the overall architecture of the system
hierarchy. However, for this Lab, you will write/modify
<strong>only</strong> `helloworld_service.go`, `driver.go`, `schedule.go`, and
`worker.go`. You will not be able to submit other files or
modules. In other words, any helper functions must reside within
these listed files.

## Implementation

### Part A: A Hello World Lambda plugin

Under directory `plugins`, we supply you a `helloworld_service`
Lambda plugin file, which is missing the basic logic in
`DoService()`. The first task you need to accomplish is to fix
that simple piece by adding logic to deserialize (unmarshal) the
argument passed in as a *raw* byte[] array. The second task is
also easy - just to print out the deserialized task number
from `DoService()`.

To build your plugin Lambda function library, make sure the
environmental variable `$GOPATH` has been correctly
configured:

```bash
$ cd $HOME
$ mkdir go
$ cd go
# Go needs $GOPATH to be set to the directory containing "cs675-spring20-labs", and the dirs "main", "plugins", and "serverless" therein.
$ export GOPATH="$PWD"
```

And then from `$HOME/go`, create a directory called `src/cs675-spring20-lab1`, 
and from there, clone your **private** repo:

```bash
$ mkdir -p src/cs675-spring20-labs
$ cd src/cs675-spring20-labs
$ git clone git@git.gmu.edu:cs675-spring20-labs/lab1.git
$ cd lab1
$ ls
README.md main/ plugins/ serverless/
```

Then build your first plugin Lambda function library with the
following commands:

```bash
$ cd plugins/
$ go build --buildmode=plugin -o helloworld_service.so helloworld_service.go
$ ls
helloworld_service.so helloworld_service.go
```


### Part B: Implement RPC-based Lambda plugin registration protocol

This part is to get you familiar with RPC distributed system
programming using Go. 

**First, the client part of the RPC protocol.**
The `serverless/driver.go` code we give you is missing one
crucial piece: the function that registers the plugin Lambda function
that you have just implemented. The plugin registration logic is
carried out by the `registerService()` function in
`serverless/driver.go`. The comments in this file should point
you in the right direction.

**Second, the worker's RPC server startup.**
Each worker has its own RPC server, which is created and started by
calling the `startRPCServer()` function (see
`main/worker.go`). You'll need to finish the implementation
of the `startRPCServer()` function at the worker side. Refer
to the implementation of the driver's RPC server to get a sense.

**Third, the worker-side RPC method functions to register and invoke
a plugin Lambda function service.**
Two method functions, missing their core logic, need to be fixed: 
`RegisterService()` and `InvokeService()`.  In
particular, `RegisterService()` uses Go's plugin feature to
dynamically load the binary of the compiled library into the worker
process' address space. Read Go's <a href="https://golang.org/pkg/plugin/">package plugin</a> and learn
how to use it in your code.
`InvokeService()` is called when the driver schedules a task
on a worker. Read the comments (hints) in these two functions, which
you may find helpful.

### Part C: Task scheduling (or dispatching :-)

Well, this part should really be called **dispatching**
rather than **scheduling**, though what you will be looking at is a
function called `schedule()` located at
`serverless/schedule.go`. In this Lab, your code just simply
needs to dispatch the task requests to the workers. In Lab 2, you
will get to handle some level of multi-phase scheduling (for
**MapReduce**).

The `schedule()` function is called from the
`driver.go` file to dispatch a configurable number (though it
is hard-coded as *10* just for now) of tasks that have already been
registered and loaded into the workers. 

We provide you with a code template that defines the Lambda plugin's
parameter `struct`, declares a `readyChan` to track the
workers that are ready to execute a task. Your task is to finish the
implementation of the `schedule()` function so that the driver
dispatches a total of `nTasks` tasks across a cluster of
workers available, in a way that is similar to a **FIFO (First-in
First-out)** or **RR (Round Robin)** policy - **whichever worker
completes its previously assigned task will get enqueued back into
the** `readyChan`. 

Specifically, `invokeService()` is where a task request gets
really sent out. Fill out the missing piece in there.  The next
missing piece to be fixed is a `for` loop that loops over all the tasks;
you may find `select` comes in handy: inside the `for`
loop, the `select` switches between two event sources:
`registerChan`, and `readyChan`.  This is also where
you get to connect the different pieces of the RPC framework together
- `registerChan` holds the workers that successfully get
registered at the driver, and that is where the task scheduler
initially starts dragging workers from.

## Deployment and Testing

You can test your implementation in a semi-distributed environment,
where the client (and thus the driver) and the worker are running as
separate processes but on the same server but communicate through
TCP-based RPC.

To deploy, first, run the client as the application which creates a
driver. The driver listens on `localhost:1234`:

```bash
$ cd main
$ go run client.go localhost:1234 helloworld_service
2020/01/20 02:39:37 rpc.Register: method "Lock" has 1 input parameters; needs exactly three
2020/01/20 02:39:37 rpc.Register: method "Run" has 2 input parameters; needs exactly three
2020/01/20 02:39:37 rpc.Register: method "Unlock" has 1 input parameters; needs exactly three
2020/01/20 02:39:37 rpc.Register: method "Wait" has 1 input parameters; needs exactly three
Driver: enter the worker registration service loop...
```

Ignore the first four lines of the output. Driver is now running and
waiting for workers to register.

On a separate shell window, run one worker (you can deploy however
many worker processes you like) that listens on
`localhost:1235` and connects to the driver located
at `localhost:1234`:

```bash
$ go run worker.go localhost:1235 localhost:1234
2020/01/20 02:44:06 rpc.Register: method "Lock" has 1 input parameters; needs exactly three
2020/01/20 02:44:06 rpc.Register: method "Unlock" has 1 input parameters; needs exactly three
Successfully registered worker localhost:1235
Successfully registered new service helloworld_service
...  # rest of the output ignored 
```


## Resources and Advice
<ul class="hints">
  <li>
    a good read on Go's concurrency patterns:
    <a href="https://blog.golang.org/pipelines">Go Blog on concurrency</a>.
  </li>
  <li>
	read the document of 
    <a href="https://golang.org/pkg/plugin/">Go plugin package</a>.
  </li>
  <li>
    read this short tutorial about <a href="https://gobyexample.com/select">Go by Example: Select</a>
    to give you a sense about how <tt>select</tt> can be
used in combination with goroutines.
  </li>
</ul>

## Point Distribution

<table>
<tr><th>Component</th><th>Points</th></tr>
<tr><td>Driver-to-Worker RPC</td><td>15</td></tr>
<tr><td>Scheduler</td><td>10</td></tr>
<tr><td>Plugin</td><td>5</td></tr>
</table>


## Submitting Lab 1

1. **Submit the electronic version**

You hand in your lab assignment exactly as you've been letting us know your progress:

```bash
$ git commit -am "[you fill me in]"
$ git tag -a -m "i finished lab 1" lab1-handin
$ git push origin master lab1-handin
```

You should verify that you are able to see your final commit and your
lab1-handin tag on the GitLab page in your repository for this lab.

We will use the timestamp of your **last** tag for the
purpose of calculating late days, and we will only grade that version of the
code. (We'll also know if you backdate the tag, don't do that.)

You will need to share your private repository with me (the instructor)
(my GitLab ID is the same as my mason email ID: `yuecheng`).

2. **Schedule a meeting and discuss**

As a second part of the submission, you'll meet with me and explain what you
did for Lab 1. Hopefully we will use the office hour for this
after the due of Lab 1. We will also have a signup sheet as the date 
approaches, and I'll also give a little more detail in class.

</p>

<h2>Acknowledgements</h2>
<p>Part of this lab is adapted from MIT's 6.824 course. Thanks to
Frans Kaashoek, Robert Morris, and Nickolai Zeldovich for their
support.</p>
