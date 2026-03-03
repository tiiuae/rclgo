# Service package(s) example
This is an example of how to use a custom ROS 2 service, to deploy a ROS 2 service node, alongside a ROS 2 client. `example_srvs` contains a ROS 2 package with the service definition that is used by both `client` and `service`.

First, you'll need to build `example_srvs`. Ensure that ROS 2 is installed and sourced, and you have the C build tools, colcon, and rosidl installed. Run the following command in the `example_srvs` directory:
```
colcon build
```
After this, source the package by running:
```
. install/setup.bash
```
Switch to the `service` directory and generate Go bindings with:
```
go generate
```
Now we can build and run the node with:
```
go build
./service
```
In another terminal, source your ROS 2 and navigate back into the `example_srvs` directory and run:
```
. install/setup.bash
```
Now navigate to the `client` directory and generate Go bindings, then build and run the program:
```
go generate
go build
./client
```
This will send a service call to our `service` node using the service `/test/add`, returning the sum of two numbers.