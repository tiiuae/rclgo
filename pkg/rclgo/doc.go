/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

/*
Deliberate trial and error have been conducted in finding the best way of
interfacing with rcl or rclc.

rclc was initially considered, but: Executor subscription callback doesn't
include the subscription, only the ros2 message. Thus we cannot intelligently
and dynamically dispatch the ros2 message to the correct subscription callback
on the golang layer. rcl wait_set has much more granular way of defining how the
received messages are handled and allows for a more Golang-way of handling
dynamic callbacks.
*/
package rclgo
