rclgo the ROS2 client library Golang wrapper
============================================

Installation
------------

    $ go get github.com/tiiuae/rclgo
    $ rclgo-gen generate

Commandline client
------------------

Mimics the official rcl-command

    rclgo topic echo /topic/name std_msgs.ColorRGBA

ROS2 message converter
----------------------

rclgo expects a Golang-implementation of all the ROS2 messages to exists.
To use rclgo with your set of ROS2 plugins and modules, you need to generate the Golang-bindings before first use.

    rclgo-gen generate /opt/ros/foxy/share/px4_msgs/msg/AdcReport.msg

Usage
-----

See the rclgo commandline client source code:

[Subscription](cmd/rclgo/topic-echo.go)
[Publisher](cmd/rclgo/topic-pub.go)
