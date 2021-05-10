/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package ros2

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/tiiuae/rclgo/pkg/datagenerator"
	"github.com/tiiuae/rclgo/pkg/ros2/ros2_type_dispatcher"
	"github.com/tiiuae/rclgo/pkg/ros2/ros2types"
)

/*
Define action bundles which generate typical use-cases with minimal effort
*/

/*
Creates a ROS2 RCL context with a single subscriber subscribing to the given topic and waiting for termination via the given/returned context.
All parameters except the first one are optional.
*/
func SubscriberBundle(ctx context.Context, rclContext *Context, wg *sync.WaitGroup, namespace, nodeName, topicName, msgTypeName string, rosArgs *RCLArgs, subscriberCallback SubscriptionCallback) (*Context, *RCLErrors) {
	var err error
	var errs *RCLErrors
	var msgType ros2types.ROS2Msg
	rclContext, wg, msgType, errs = bundleDefaults(rclContext, wg, &namespace, &nodeName, &topicName, &msgTypeName, rosArgs)
	if errs != nil {
		return rclContext, errs
	}

	rclNode, err := rclContext.NewNode(nodeName, namespace)
	if err != nil {
		return rclContext, RCLErrorsPut(errs, err)
	}

	ros2msgClone := msgType.Clone()
	subscription, err := rclNode.NewSubscription(topicName, ros2msgClone, subscriberCallback)
	if err != nil {
		return rclContext, RCLErrorsPut(errs, err)
	}

	waitSet, err := rclContext.NewWaitSet(1000 * time.Millisecond)
	if err != nil {
		return rclContext, RCLErrorsPut(errs, err)
	}
	waitSet.AddSubscriptions(subscription)

	waitSet.RunGoroutine(ctx)

	return rclContext, errs
}

func PublisherBundle(rclContext *Context, wg *sync.WaitGroup, namespace, nodeName, topicName, msgTypeName string, rosArgs *RCLArgs) (*Context, *Publisher, *RCLErrors) {
	var err error
	var errs *RCLErrors
	var msgType ros2types.ROS2Msg
	rclContext, _, msgType, errs = bundleDefaults(rclContext, wg, &namespace, &nodeName, &topicName, &msgTypeName, rosArgs)
	if errs != nil {
		return rclContext, nil, errs
	}

	rclNode, err := rclContext.NewNode(nodeName, namespace)
	if err != nil {
		return rclContext, nil, RCLErrorsPut(errs, err)
	}

	publisher, err := rclNode.NewPublisher(topicName, msgType)
	if err != nil {
		return rclContext, nil, RCLErrorsPut(errs, err)
	}

	return rclContext, publisher, errs
}

func PublisherBundleTimer(ctx context.Context, rclContext *Context, wg *sync.WaitGroup, namespace, nodeName, topicName, msgTypeName string, rosArgs *RCLArgs, interval time.Duration, payload string, publisherCallback func(*Publisher, ros2types.ROS2Msg) bool) (*Context, *RCLErrors) {
	var errs *RCLErrors
	var publisher *Publisher
	rclContext, publisher, errs = PublisherBundle(rclContext, wg, namespace, nodeName, topicName, msgTypeName, rosArgs)
	if errs != nil {
		return rclContext, errs
	}

	timer, err := rclContext.NewTimer(interval, func(timer *Timer) {
		// It would be smarter to allocate memory for the ros2msg outside the timer callback, but this way the tests can test for memory leaks too using this same codebase.
		ros2msg, err_yaml := ros2_type_dispatcher.TranslateMsgPayloadYAMLToROS2Msg(strings.ReplaceAll(payload, "\\n", "\n"), publisher.Ros2MsgType)
		if err_yaml != nil {
			errs = RCLErrorsPut(errs, ErrorsCastC(1003, fmt.Sprintf("Error '%v' unmarshalling YAML '%s' to ROS2 message type '%s'", err_yaml, payload, msgTypeName)))
		}
		if publisherCallback != nil {
			if publisherCallback(publisher, ros2msg) {
				publisher.Publish(ros2msg)
			}
		} else {
			publisher.Publish(ros2msg)
		}
	})
	if err != nil {
		return rclContext, RCLErrorsPut(errs, err)
	}

	waitSet, err := rclContext.NewWaitSet(1000 * time.Millisecond)
	if err != nil {
		return rclContext, RCLErrorsPut(errs, err)
	}
	waitSet.AddTimers(timer)

	waitSet.RunGoroutine(ctx)

	return rclContext, errs
}

/*
bundleDefaults creates a default context from the given parameters.
*/
func bundleDefaults(rclContext *Context, wg *sync.WaitGroup, namespace, nodeName, topicName, msgTypeName *string, rosArgs *RCLArgs) (*Context, *sync.WaitGroup, ros2types.ROS2Msg, *RCLErrors) {
	var errs *RCLErrors
	var err error

	if rosArgs == nil {
		rosArgs, err = NewRCLArgs("")
		if err != nil {
			errs = RCLErrorsPut(errs, err)
		} else {
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()
			os.Args = []string{}
			rosArgs, err = NewRCLArgs("")
			if err != nil {
				return nil, wg, nil, RCLErrorsPut(errs, err)
			}
		}
	}

	if rclContext == nil {
		rclContext, err = NewContext(wg, 0, rosArgs)
		if err != nil {
			return rclContext, wg, nil, RCLErrorsPut(errs, err)
		}
	} else {
		if wg == nil && rclContext.WG != nil {
			// wg already exists in the RCL context
		} else if wg != nil {
			rclContext.WG = wg
		} else {
			rclContext.WG = &sync.WaitGroup{}
		}
	}

	if *nodeName == "" {
		*nodeName = datagenerator.NodeName()
	}

	ros2msg, ok := ros2_type_dispatcher.TranslateROS2MsgTypeNameToType(*msgTypeName)
	if !ok {
		return rclContext, wg, ros2msg, RCLErrorsPut(errs, ErrorsCastC(1003, fmt.Sprintf("No ROS2 Message mapping from type '%s'", *msgTypeName)))
	}
	return rclContext, wg, ros2msg, errs
}
