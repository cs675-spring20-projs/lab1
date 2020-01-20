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
The serverless package (located at <tt>src/serverless</tt>) provides
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
README.md src
# Go needs $GOPATH to be set to the directory containing "src"
$ export GOPATH="$PWD"
$ cd src
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
