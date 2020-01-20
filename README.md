# CS 675 Lab 1: Serverless and RPC

<h2>Introduction</h2>
<p>
In Lab 1, you will build a client and server framework as a way to
to learn the Go programming language and as a way to learn about RPCs
in distributed systems. 
</p>

<h2>Software</h2>
<p>
You'll implement this Lab (and all the Labs) in  <a
href="http://www.golang.org/">Go</a>. The Go website contains lots
of tutorial information which you may want to look at.
</p>

<p>
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
</p>

<h3>Getting familiar with the source</h3>
<p>
The serverless package (located at <tt>serverless</tt>) provides
a simple serverless library with a partially implemented RPC plugin
framework. Applications would normally call <tt>Run()</tt> located in
<tt>serverless/driver.go</tt> to register a Lambda function (i.e., a 
Go plugin service library) and start a configurable number of tasks
(executing the registered Lambda function).
</p>

<p>
The flow of the RPC client-server implementation is as follows:
<ol>
	<li>
		The user (i.e., you) provides one or multiple Lambda functions in
the form of a Go plugin library (see the skeleton code
<tt>helloworld_service.go</tt> provided under <tt>plugins/</tt>).
	</li>
	<li>
		A driver is created with this knowledge when running
<tt>client</tt> (see <tt>main/client.go</tt>). It spins up an RPC
server (see <tt>serverless/driver.go</tt>), and waits for workers to
register (using the RPC call <tt>Register()</tt> defined in
<tt>serverless/driver.go</tt>). 
	</li>
	<li>
		One or multiple worker processes are created when running
<tt>worker</tt> (see <tt>main/worker.go</tt>). Each worker spins up
an RPC server, and registers itself at the driver, and waits for the
driver to register Lambda functions and schedule tasks.
	</li>
	<li>
		The driver registers the Lambda function specified as
<tt>serviceName</tt> (the second command-line argument when running
<tt>driver.go</tt>) by calling <tt>prepareService()</tt> (see
<tt>serverless/driver.go</tt>), which issues the
<tt>RegisterService()</tt> RPC call on a worker.  The driver then
schedules the tasks by calling <tt>schedule()</tt> (see
<tt>serverless/schedule.go</tt>). In <tt>schedule()</tt>, the driver
issues the <tt>InvokeService()</tt> RPC call to execute a
plugin Lambda function on a worker.
	</li>
	<li>
		The driver sends a <tt>Shutdown()</tt> RPC to each of its
workers, and then shuts down its own RPC server.
	</li>
</ol>

You should look through the files in the whole framework
implementation, as reading them might be useful to understand how the
other methods fit into the overall architecture of the system
hierarchy. However, for this Lab, you will write/modify
<strong>only</strong> <tt>driver.go</tt>, <tt>schedule.go</tt>, and
<tt>worker.go</tt>. You will not be able to submit other files or
modules. In other words, any helper functions must reside within
these listed files.
</p> 

<h2>Implementation</h2>

<h3>Part A: A Hello World Lambda plugin</h3>

<p>
Under <tt>plugins</tt>, we supply you a <tt>helloworld_service</tt>
Lambda plugin file, which is missing the basic logic in
<tt>DoService()</tt>. The first task you need to accomplish is to fix
that simple piece by adding logic to deserialize (unmarshal) the
argument passed in as a <tt>raw</tt> byte[] array. The second task is
also easy peasy -- just to print out the deserialized task number
from <tt>DoService()</tt>.
</p>

<p>
To build your plugin Lambda function library, make sure the
environmental variable <tt>$GOPATH</tt> has been correctly
configured:

<pre>
$ cd cs675-spring20-labs
$ ls
README.md main/ plugins/ serverless/
# Go needs $GOPATH to be set to the directory containing "main", "plugins", and "serverless".
$ export GOPATH="$PWD"
</pre>

<p>
Then build your first plugin Lambda function library with the
following commands:
<pre>
$ cd plugins/
$ go build --buildmode=plugin -o helloworld_service.so helloworld_service.go
$ ls
helloworld_service.so helloworld_service.go
</pre>

</p>

<h3>Part B: Implement RPC-based Lambda plugin registration protocol</h3>

<p>
This part is to get you familiar with RPC distributed system
programming using Go. 

<p>
First, the client part of the RPC protocol.
The <tt>serverless/driver.go</tt> code we give you is missing one
crucial piece: the function that register the plugin Lambda function
that you have just implemented. The plugin registration logic is
carried out by the <tt>registerService()</tt> function in
<tt>serverless/driver.go</tt>. The comments in this file should point
you in the right direction.

<p>
Second, the worker's RPC server startup.
Each worker has its own RPC server, which is created and started by
calling the <tt>startRPCServer()</tt> function (see
<tt>main/worker.go</tt>). You will need to finish the implementation
of the <tt>startRPCServer()</tt> function at the worker side. Refer
to the implementation of the driver's RPC server to get a sense.


<p>
Third, the worker-side RPC functions to register and invoke a plugin Lambda
function service.
Two functions, missing their core logic, need to be fixed: 
<tt>RegisterService()</tt> and <tt>InvokeService()</tt>.  In
particular, <tt>RegisterService()</tt> uses Go's plugin feature to
dynamically load the compiled library binary into the worker's
address space. Read Go's <a
href="https://golang.org/pkg/plugin/">package plugin</a> and learn
how to use it in your code.
<tt>InvokeService()</tt> is called when the driver schedules a task
on the worker. You may find the comments in these two functions 
helpful.
</p>

<h3>Part C: Task scheduling (or dispatching :-)</h3>

<p>
Well, this part should really be called <strong>dispatching</strong>
rather than *scheduling*, though what you will be looking at is a
function called <tt>schedule()</tt> located at
<tt>serverless/schedule.go</tt>. In this Lab, your code just simply
needs to dispatch the task requests to the workers. In Lab 2, you
will get to handle some level of phase-based scheduling.

<p>
The <tt>schedule()</tt> function is called from the
<tt>driver.go</tt> file to dispatch a configurable number (though it
is hard-coded as 10) of tasks that have already been registered and
loaded into the workers. 

We provide you with a code template that defines the Lambda plugin's
parameter struct, declares a <tt>readyChan</tt> to track the
workers that are ready to execute a task. Your task is to finish the
implementation of the <tt>schedule()</tt> function so that the driver
dispatches a total of <tt>nTasks</tt> tasks across a cluster of
workers available, in a FIFO (First-in First-out) manner -- whichever
worker completes its previously assigned task will get enqueued back
into the <tt>readyChan</tt>. 

Specifically, <tt>invokeService()</tt> is where the task dispatching
really happens. Fill out the missing piece in there.  The next
missing piece to be fixed is a <tt>for</tt> loop that loops over all the tasks;
you may find <tt>select</tt> comes in handy: inside the <tt>for</tt>
loop, the <tt>select</tt> switches between two event sources:
<tt>registerChan</tt>, and <tt>readyChan</tt>.  This is also where
you get to connect the different pieces of the RPC framework together
-- <tt>registerChan</tt> holds the workers that successfully get
registered at the driver, and that is where the task scheduler
initially starts dragging workers from.

<h2>Deployment and Testing</h2>

<p>
You can test your implementation in a semi-distributed environment,
where the client (and thus the driver) and the worker are running as
separate processes but on the same server but communicate through
TCP-based RPC.

To deploy, first, run the client as the application which creates a
driver:

<pre>
$ cd main
$ go run client.go localhost:1234 helloworld_service
2020/01/20 02:39:37 rpc.Register: method "Lock" has 1 input parameters; needs exactly three
2020/01/20 02:39:37 rpc.Register: method "Run" has 2 input parameters; needs exactly three
2020/01/20 02:39:37 rpc.Register: method "Unlock" has 1 input parameters; needs exactly three
2020/01/20 02:39:37 rpc.Register: method "Wait" has 1 input parameters; needs exactly three
Driver: enter the worker registration service loop...
</pre>

Ignore the first four lines of the output. Driver is now running and
waiting for workers to register.

On a separate shell window, run one worker (you can deploy however
many worker processes you like) that listens on
<tt>localhost:1235</tt> and connects to the driver located
at <tt>localhost:1234</tt>:

<pre>
$ go run worker.go localhost:1235 localhost:1234
2020/01/20 02:44:06 rpc.Register: method "Lock" has 1 input parameters; needs exactly three
2020/01/20 02:44:06 rpc.Register: method "Unlock" has 1 input parameters; needs exactly three
Successfully registered worker localhost:1235
Successfully registered new service helloworld_service
...  # rest of the output ignored 
</pre>

</p>


<h2>Resources and Advice</h2>
<ul class="hints">
  <li>
    a good read on what strings are in Go is the
    <a href="http://blog.golang.org/strings">Go Blog on strings</a>.
  </li>
  <li>
	read the document of 
    <a href="https://golang.org/pkg/plugin/"><tt>Go plugin package</tt></a>.
  </li>
  <li>
    read the <a href="https://gobyexample.com/select">Go by Example: Select</a>
    is give you a sense about how <tt>select</tt> can be
used in combination with goroutines.
  </li>
</ul>

## Point Distribution

TBD

## Submitting Lab 1

TBD

We will use the timestamp of your **last** tag for the
purpose of calculating late days, and we will only grade that version of the
code. (We'll also know if you backdate the tag, don't do that.)

</p>
