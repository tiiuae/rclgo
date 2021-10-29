/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package rclgo

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/tiiuae/rclgo/pkg/datagenerator"
	"github.com/tiiuae/rclgo/pkg/rclgo/typemap"
	"github.com/tiiuae/rclgo/pkg/rclgo/types"
	"gopkg.in/yaml.v2"
)

/*
Define action bundles which generate typical use-cases with minimal effort
*/

/*
Creates a ROS2 RCL context with a single subscriber subscribing to the given topic and waiting for termination via the given/returned context.
All parameters except the first one are optional.
*/
func SubscriberBundle(ctx context.Context, rclContext *Context, wg *sync.WaitGroup, namespace, nodeName, topicName, msgTypeName string, rosArgs *Args, subscriberCallback SubscriptionCallback) (*Context, error) {
	c, _, err := SubscriberBundleReturnWaitSet(ctx, rclContext, wg, namespace, nodeName, topicName, msgTypeName, rosArgs, subscriberCallback)
	return c, err
}

func SubscriberBundleReturnWaitSet(ctx context.Context, rclContext *Context, wg *sync.WaitGroup, namespace, nodeName, topicName, msgTypeName string, rosArgs *Args, subscriberCallback SubscriptionCallback) (*Context, *WaitSet, error) {
	var errs error
	var msgType types.MessageTypeSupport
	rclContext, wg, msgType, errs = bundleDefaults(rclContext, wg, &namespace, &nodeName, &topicName, &msgTypeName, rosArgs)
	if errs != nil {
		return rclContext, nil, errs
	}

	rclNode, err := rclContext.NewNode(nodeName, namespace)
	if err != nil {
		return rclContext, nil, multierror.Append(errs, err)
	}

	subscription, err := rclNode.NewSubscription(topicName, msgType, subscriberCallback)
	if err != nil {
		return rclContext, nil, multierror.Append(errs, err)
	}

	waitSet, err := rclContext.NewWaitSet(1000 * time.Millisecond)
	if err != nil {
		return rclContext, nil, multierror.Append(errs, err)
	}
	waitSet.AddSubscriptions(subscription)

	waitSet.RunGoroutine(ctx)

	return rclContext, waitSet, errs
}

func PublisherBundle(rclContext *Context, wg *sync.WaitGroup, namespace, nodeName, topicName, msgTypeName string, rosArgs *Args) (*Context, *Publisher, error) {
	var errs error
	var msgType types.MessageTypeSupport
	rclContext, _, msgType, errs = bundleDefaults(rclContext, wg, &namespace, &nodeName, &topicName, &msgTypeName, rosArgs)
	if errs != nil {
		return rclContext, nil, errs
	}

	rclNode, err := rclContext.NewNode(nodeName, namespace)
	if err != nil {
		return rclContext, nil, multierror.Append(errs, err)
	}

	publisher, err := rclNode.NewPublisher(topicName, msgType, nil)
	if err != nil {
		return rclContext, nil, multierror.Append(errs, err)
	}

	return rclContext, publisher, errs
}

func PublisherBundleTimer(ctx context.Context, rclContext *Context, wg *sync.WaitGroup, namespace, nodeName, topicName, msgTypeName string, rosArgs *Args, interval time.Duration, payload string, publisherCallback func(*Publisher, types.Message) bool) (*Context, error) {
	var errs error
	var publisher *Publisher
	rclContext, publisher, errs = PublisherBundle(rclContext, wg, namespace, nodeName, topicName, msgTypeName, rosArgs)
	if errs != nil {
		return rclContext, errs
	}

	timer, err := rclContext.NewTimer(interval, func(timer *Timer) {
		// It would be smarter to allocate memory for the ros2msg outside the timer callback, but this way the tests can test for memory leaks too using this same codebase.
		ros2msg := publisher.typeSupport.New()
		err_yaml := yaml.Unmarshal([]byte(strings.ReplaceAll(payload, "\\n", "\n")), ros2msg)
		if err_yaml != nil {
			errs = multierror.Append(errs, errorsCastC(1003, fmt.Sprintf("Error '%v' unmarshalling YAML '%s' to ROS2 message type '%s'", err_yaml, payload, msgTypeName)))
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
		return rclContext, multierror.Append(errs, err)
	}

	waitSet, err := rclContext.NewWaitSet(1000 * time.Millisecond)
	if err != nil {
		return rclContext, multierror.Append(errs, err)
	}
	waitSet.AddTimers(timer)

	waitSet.RunGoroutine(ctx)

	return rclContext, errs
}

/*
bundleDefaults creates a default context from the given parameters.
*/
func bundleDefaults(rclContext *Context, wg *sync.WaitGroup, namespace, nodeName, topicName, msgTypeName *string, rosArgs *Args) (*Context, *sync.WaitGroup, types.MessageTypeSupport, error) {
	var err, errs error

	if rosArgs == nil {
		rosArgs, _, err = ParseArgs(os.Args)
		if err != nil {
			return nil, wg, nil, err
		}
	}

	if rclContext == nil {
		rclContext, err = NewContext(wg, 0, rosArgs)
		if err != nil {
			return rclContext, wg, nil, multierror.Append(errs, err)
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

	ros2msg, ok := typemap.GetMessage(*msgTypeName)
	if !ok {
		return rclContext, wg, ros2msg, multierror.Append(errs, errorsCastC(1003, fmt.Sprintf("No ROS2 Message mapping from type '%s'", *msgTypeName)))
	}
	return rclContext, wg, ros2msg, errs
}
