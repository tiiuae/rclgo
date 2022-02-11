package rclgo_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-multierror"
	. "github.com/smartystreets/goconvey/convey"
	std_msgs "github.com/tiiuae/rclgo/internal/msgs/std_msgs/msg"
	test_msgs "github.com/tiiuae/rclgo/internal/msgs/test_msgs/msg"
	"github.com/tiiuae/rclgo/pkg/rclgo"
	"github.com/tiiuae/rclgo/pkg/rclgo/types"
	"gopkg.in/yaml.v3"
)

func TestPubSub(t *testing.T) {
	var rclContextPub *rclgo.Context
	var rclContextSub *rclgo.Context
	var publisher *rclgo.Publisher
	subChan := make(chan *std_msgs.ColorRGBA, 1)
	subErrChan := make(chan error, 1)
	subCtx, cancelSubCtx := context.WithCancel(context.Background())
	defer cancelSubCtx()
	Convey(
		"Scenario: Publisher publishes, Subscriber subscribes, garbage is collected",
		t,
		func() {
			Convey("Given a Subscriber", func() {
				var err error
				rclContextSub, err = newContextWithSubscriber(
					"sub",
					"/test",
					"/topic",
					std_msgs.ColorRGBATypeSupport,
					func(s *rclgo.Subscription) {
						buf, _, err := s.TakeSerializedMessage()
						if err != nil {
							panic(fmt.Sprint("failed to take message: ", err))
						}
						msg, err := rclgo.Deserialize(buf, std_msgs.ColorRGBATypeSupport)
						if err != nil {
							panic(fmt.Sprint("failed to deserialize message: ", err))
						}
						newBuf, err := rclgo.Serialize(msg)
						if err != nil {
							panic(fmt.Sprint("failed to reserialize message: ", err))
						} else if !bytes.Equal(buf, newBuf) {
							panic(fmt.Sprintf("reserialized message (%#v) is different from the original (%#v)", newBuf, buf))
						}
						subChan <- msg.(*std_msgs.ColorRGBA)
					},
				)
				So(err, ShouldBeNil)
				go func() { subErrChan <- rclContextSub.Spin(subCtx) }()
				time.Sleep(200 * time.Millisecond)
			})
			Convey("And a Publisher", func() {
				var err error
				rclContextPub, publisher, err = newContextWithPublisher(
					nil,
					"pub",
					"/test",
					"/topic",
					std_msgs.ColorRGBATypeSupport,
				)
				So(err, ShouldBeNil)
			})
			Convey("When the Publisher publishes", func() {
				err := publishColorRGBA(publisher, 1.55, 2.66, 3.77, 4.88)
				So(err, ShouldBeNil)
			})
			Convey("Then the Subscriber receives", func() {
				receiveColorRGBA(subChan, 1.55, 2.66, 3.77, 4.88)
			})
			Convey("When the Publisher publishes again", func() {
				err := publishColorRGBA(publisher, 0.00, 1.00, 2.00, 3.00)
				So(err, ShouldBeNil)
			})
			Convey("Then the Subscriber receives again", func() {
				receiveColorRGBA(subChan, 0.00, 1.00, 2.00, 3.00)
			})
			Convey("When the Subscriber context is canceled", func() {
				var err error
				cancelSubCtx()
				timeOut(
					2000,
					func() { err = <-subErrChan },
					"Subscriber waitGroup waiting to finish",
				)
				So(err, shouldContainError, context.Canceled)
			})
			Convey("And the Subscriber context is GC'd", func() {
				errs := rclContextSub.Close()
				So(errs, ShouldBeNil)
			})
			Convey(
				"And the Publisher publishes to a Topic with no Subscribers",
				func() {
					err := publishColorRGBA(publisher, 0.00, 1.00, 2.00, 3.00)
					So(err, ShouldBeNil)
				},
			)
			Convey("Then the Subscriber cannot receive the message", func() {
				So(len(subChan), ShouldEqual, 0)
			})
			Convey("And the Publisher context is GC'd", func() {
				errs := rclContextPub.Close()
				So(errs, ShouldBeNil)
			})
		},
	)
}

func TestMultipleSubscribersInSingleWaitSet(t *testing.T) {
	var (
		rclCtxPub, rclCtxSub *rclgo.Context
		pub1, pub2           *rclgo.Publisher
		topic1Chan           = make(chan receiveResult, 1)
		topic2Chan           = make(chan receiveResult, 1)
		errChan              = make(chan error, 2)
	)
	defer func() {
		if rclCtxPub != nil {
			rclCtxPub.Close()
		}
		if rclCtxSub != nil {
			rclCtxSub.Close()
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	SetDefaultFailureMode(FailureContinues)
	Convey(
		"Scenario: Publishers publishe, multiple Subscribers listen",
		t,
		func() {
			Convey("Given RCL contexts and waitset", func() {
				var err error
				rclCtxSub, err = newDefaultRCLContext()
				So(err, ShouldBeNil)
				rclCtxPub, err = newDefaultRCLContext()
				So(err, ShouldBeNil)
			})
			Convey("And a subscriber on the first topic", func() {
				node, err := rclCtxSub.NewNode("sub1", "/test")
				So(err, ShouldBeNil)
				_, err = node.NewSubscription(
					"/topic1",
					std_msgs.StringTypeSupport,
					sendToChan(topic1Chan),
				)
				So(err, ShouldBeNil)
			})
			Convey("And a subscriber on the second topic", func() {
				node, err := rclCtxPub.NewNode("sub2", "/test")
				So(err, ShouldBeNil)
				_, err = node.NewSubscription(
					"/topic2",
					std_msgs.StringTypeSupport,
					sendToChan(topic2Chan),
				)
				So(err, ShouldBeNil)
			})
			Convey("And a publisher on the first topic", func() {
				node, err := rclCtxSub.NewNode("pub1", "/test")
				So(err, ShouldBeNil)
				pub1, err = node.NewPublisher(
					"/topic1",
					std_msgs.StringTypeSupport,
					nil,
				)
				So(err, ShouldBeNil)
			})
			Convey("And a publisher on the second topic", func() {
				node, err := rclCtxPub.NewNode("pub2", "/test")
				So(err, ShouldBeNil)
				pub2, err = node.NewPublisher(
					"/topic1",
					std_msgs.StringTypeSupport,
					nil,
				)
				So(err, ShouldBeNil)
			})
			Convey("And the waitset is started", func() {
				go func() { errChan <- rclCtxSub.Spin(ctx) }()
				go func() { errChan <- rclCtxPub.Spin(ctx) }()
			})
			Convey("When the first pub312lisher publishes", func() {
				publishString(pub1, "message on topic 1")
			})
			Convey("Then the first subscriber receives", func() {
				receiveString(topic1Chan, "message on topic 1")
			})
			Convey("But the second subscriber doesn't receive", func() {
				receiveNothing(topic2Chan)
			})
			Convey("When the second publisher publishes", func() {
				publishString(pub2, "message on topic 2")
			})
			Convey("Then the first subscriber doesn't receive", func() {
				receiveNothing(topic2Chan)
			})
			Convey("But the second subscriber receives", func() {
				receiveString(topic1Chan, "message on topic 2")
			})
			Convey("When the context is canceled", func() {
				cancel()
			})
			Convey("Then all subscribers stop", func() {
				var err error
				timeOut(2000, func() {
					err = multierror.Append(<-errChan, <-errChan)
				}, "Subscriber waitGroup waiting to finish")
				So(err, shouldContainError, context.Canceled)
			})
			Convey("And RCL contexts are freed", func() {
				So(rclCtxSub.Close(), ShouldBeNil)
				So(rclCtxPub.Close(), ShouldBeNil)
			})
		},
	)
}

func TestMultipleTimersInSingleWaitSet(t *testing.T) {
	var (
		rclCtx     *rclgo.Context
		timer1Chan = make(chan struct{}, 1)
		timer2Chan = make(chan struct{}, 1)
		errChan    = make(chan error, 1)
	)
	defer func() {
		if rclCtx != nil {
			rclCtx.Close()
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Convey("Scenario: Multiple timers on single wait set", t, func() {
		Convey("Given and RCL context and two timers", func() {
			var err error
			rclCtx, err = newDefaultRCLContext()
			So(err, ShouldBeNil)
			_, err = rclCtx.NewTimer(time.Second, func(t *rclgo.Timer) {
				timer1Chan <- struct{}{}
			})
			So(err, ShouldBeNil)
			_, err = rclCtx.NewTimer(time.Hour, func(t *rclgo.Timer) {
				timer2Chan <- struct{}{}
			})
			So(err, ShouldBeNil)
		})
		Convey("When WaitSet is started", func() {
			go func() { errChan <- rclCtx.Spin(ctx) }()
		})
		Convey("Then callback of timer 1 is called", func() {
			timeOut(
				1050,
				func() { <-timer1Chan },
				"Waiting for timer 1 callback",
			)
		})
		Convey("But callback of timer 2 is not called", func() {
			receiveNothing(timer2Chan)
		})
		Convey("When the context is canceled", func() {
			cancel()
		})
		Convey("Then all subscribers stop", func() {
			var err error
			timeOut(
				2000,
				func() { err = <-errChan },
				"Subscriber waitGroup waiting to finish",
			)
			So(err, shouldContainError, context.Canceled)
		})
		Convey("And RCL contexts are freed", func() {
			So(rclCtx.Close(), ShouldBeNil)
		})
	})
}

func BenchmarkPubsubMemoryLeakAllocateInLoop(t *testing.B) {
	var messagesReceived int = 0
	fmt.Printf(
		"Mem from pmap(1) '%skB' messages '%d'\n",
		getMemReading(),
		messagesReceived,
	)
	for {
		runCtx, stopRun := context.WithCancel(context.Background())
		errChan := make(chan error, 2)
		rclContextSub, err := newContextWithSubscriber(
			"sub",
			"/test",
			"/topic",
			test_msgs.UnboundedSequencesTypeSupport,
			func(s *rclgo.Subscription) {
				var m test_msgs.UnboundedSequences
				_, err := s.TakeMessage(&m)
				if err != nil {
					fmt.Println("failed to take message:", err)
				}
				// fmt.Printf("%+v\n", c)
				messagesReceived++
			},
		)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
		go func() { errChan <- rclContextSub.Spin(runCtx) }()

		rclContextPub, errs := newContextWithPublisherTimer(
			nil,
			"pub",
			"/test",
			"/topic",
			test_msgs.UnboundedSequencesTypeSupport,
			1*time.Millisecond,
			"",
		)
		if errs != nil {
			fmt.Println("error:", errs)
			continue
		}
		go func() { errChan <- rclContextSub.Spin(runCtx) }()

		time.Sleep(1000 * time.Millisecond)

		stopRun()
		if err = <-errChan; err != nil {
			fmt.Println("error:", err)
			continue
		} else if err = <-errChan; err != nil {
			fmt.Println("error:", err)
			continue
		}
		errs = rclContextSub.Close()
		if errs != nil {
			fmt.Println("error:", errs)
		}
		errs = rclContextPub.Close()
		if errs != nil {
			fmt.Println("error:", errs)
		}
		runtime.GC()
		fmt.Printf(
			"Mem from pmap(1) '%skB' messages '%d'\n",
			getMemReading(),
			messagesReceived,
		)
	}
}

func BenchmarkPubsubMemoryLeakAllocateOutOfLoop(t *testing.B) {
	var messagesReceived int64 = 0
	fmt.Printf(
		"Mem from pmap(1) '%skB' messages '%d'\n",
		getMemReading(),
		messagesReceived,
	)
	errChan := make(chan error, 2)
	runCtx, stopRun := context.WithCancel(context.Background())
	rclContext, errs := newContextWithSubscriber(
		"sub",
		"/test",
		"/topic",
		test_msgs.UnboundedSequencesTypeSupport,
		func(s *rclgo.Subscription) {
			var m test_msgs.UnboundedSequences
			_, err := s.TakeMessage(&m)
			if err != nil {
				fmt.Println("failed to take message:", err)
			}
			// fmt.Printf("%+v\n", c)
			atomic.AddInt64(&messagesReceived, 1)
		},
	)
	if errs != nil {
		panic(errs)
	}
	go func() { errChan <- rclContext.Spin(runCtx) }()

	rclContext, errs = newContextWithPublisherTimer(
		rclContext,
		"pub",
		"/test",
		"/topic",
		test_msgs.UnboundedSequencesTypeSupport,
		1*time.Millisecond,
		"",
	)
	if errs != nil {
		panic(errs)
	}
	defer rclContext.Close()
	go func() { errChan <- rclContext.Spin(runCtx) }()
	defer stopRun()

	for {
		select {
		case err := <-errChan:
			fmt.Println("Spin returned an error:", err)
			return
		default:
			time.Sleep(1000 * time.Millisecond)
			runtime.GC()
			fmt.Printf(
				"Mem from pmap(1) '%skB' messages '%d'\n",
				getMemReading(),
				messagesReceived,
			)
		}
	}
}

func getMemReading() string {
	cmd := `pmap ` + fmt.Sprint(
		os.Getpid(),
	) + ` | tail -n 1 | grep -Po '\d+'` //  total          2102728K => 2102728
	output, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return fmt.Sprintf("Failed to execute command: %s", cmd)
	}
	return strings.TrimSpace(string(output))
}

func publishColorRGBA(p *rclgo.Publisher, r, g, b, a float32) error {
	m := std_msgs.NewColorRGBA()
	m.R = r
	m.G = g
	m.B = b
	m.A = a
	return p.Publish(m)
}

func receiveColorRGBA(subChan chan *std_msgs.ColorRGBA, r, g, b, a float32) {
	var m *std_msgs.ColorRGBA
	timeOut(1000, func() {
		m = <-subChan
	}, "Subscriber waiting for messages")
	So(m, ShouldResemble, &std_msgs.ColorRGBA{R: r, G: g, B: b, A: a})
}

func timeOut(timeoutMs int, f func(), testDescription string) bool {
	done := make(chan bool)
	go func() {
		f()
		done <- true
	}()

	select {
	case <-time.After(time.Duration(timeoutMs) * time.Millisecond):
		So("Test timeoutted!", ShouldEqual, testDescription)
		return false
	case <-done:
		So(testDescription, ShouldEqual, testDescription)
		return true
	}
}

func publishString(pub *rclgo.Publisher, s string) {
	msg := std_msgs.NewString()
	msg.Data = s
	So(pub.Publish(msg), ShouldBeNil)
}

func receiveString(subs <-chan receiveResult, expected string) {
	var m receiveResult
	timeOut(1000, func() { m = <-subs }, "Subscriber waiting for messages")
	So(m.msg, ShouldNotBeNil)
	So(m.rmi, ShouldNotBeNil)
	So(m.err, ShouldBeNil)
	So(string(m.msg.Data), ShouldEqual, expected)
}

func receiveNothing(subs interface{}) {
	i, _, _ := reflect.Select([]reflect.SelectCase{
		{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(time.After(time.Second)),
		},
		{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(subs),
		},
	})
	So(i, ShouldEqual, 0)
}

func newDefaultRCLContext() (*rclgo.Context, error) {
	return rclgo.NewContext(
		0,
		parseArgsMust("--ros-args", "--log-level", "DEBUG"),
	)
}

func newDefaultRCLContextWithOpts(opts *rclgo.ContextOptions) (*rclgo.Context, error) {
	return rclgo.NewContextWithOpts(
		parseArgsMust("--ros-args", "--log-level", "DEBUG"),
		opts,
	)
}

func parseArgsMust(args ...string) *rclgo.Args {
	a, _, err := rclgo.ParseArgs(args)
	if err != nil {
		panic("failed to parse args: " + err.Error())
	}
	return a
}

var reliableQos = func() rclgo.RmwQosProfile {
	qos := rclgo.NewRmwQosProfileDefault()
	qos.History = rclgo.RmwQosHistoryPolicyKeepAll
	qos.Durability = rclgo.RmwQosDurabilityPolicyTransientLocal
	qos.Reliability = rclgo.RmwQosReliabilityPolicyReliable
	return qos
}()

func newContextWithSubscriber(
	node, ns, topic string,
	ts types.MessageTypeSupport,
	cb rclgo.SubscriptionCallback,
) (c *rclgo.Context, err error) {
	c, err = newDefaultRCLContext()
	if err != nil {
		return nil, err
	}
	defer onErr(&err, c.Close)
	n, err := c.NewNode(node, ns)
	if err != nil {
		return nil, err
	}
	opts := rclgo.NewDefaultSubscriptionOptions()
	opts.Qos = reliableQos
	_, err = n.NewSubscriptionWithOpts(topic, ts, opts, cb)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func newContextWithPublisher(
	c *rclgo.Context,
	node, ns, topic string,
	ts types.MessageTypeSupport,
) (_ *rclgo.Context, p *rclgo.Publisher, err error) {
	if c == nil {
		c, err = newDefaultRCLContext()
		if err != nil {
			return nil, nil, err
		}
		defer onErr(&err, c.Close)
	}
	n, err := c.NewNode(node, ns)
	if err != nil {
		return nil, nil, err
	}
	defer onErr(&err, n.Close)
	opts := rclgo.NewDefaultPublisherOptions()
	opts.Qos = reliableQos
	p, err = n.NewPublisher(topic, ts, opts)
	if err != nil {
		return nil, nil, err
	}
	return c, p, nil
}

func newContextWithPublisherTimer(
	c *rclgo.Context,
	node, ns, topic string,
	ts types.MessageTypeSupport,
	interval time.Duration,
	payload string,
) (_ *rclgo.Context, err error) {
	c, p, err := newContextWithPublisher(c, node, ns, topic, ts)
	if err != nil {
		return nil, err
	}
	defer onErr(&err, c.Close)
	_, err = c.NewTimer(interval, func(t *rclgo.Timer) {
		// It would be smarter to allocate memory for the ros2msg outside the
		// timer callback, but this way the tests can test for memory leaks too
		// using this same codebase.
		ros2msg := ts.New()
		err_yaml := yaml.Unmarshal(
			[]byte(strings.ReplaceAll(payload, "\\n", "\n")),
			ros2msg,
		)
		if err_yaml == nil {
			p.Publish(ros2msg)
		} else {
			fmt.Println(rclgo.Testing_errorsCastC(1003, fmt.Sprintf("Error '%v' unmarshalling YAML '%s' to ROS2 message type '%T'", err_yaml, payload, ros2msg)))
		}
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}

func sendToChan(c chan<- receiveResult) func(s *rclgo.Subscription) {
	return func(s *rclgo.Subscription) {
		var res receiveResult
		res.rmi, res.err = s.TakeMessage(&res.msg)
		c <- res
	}
}

type receiveResult struct {
	msg std_msgs.String
	rmi *rclgo.RmwMessageInfo
	err error
}

func shouldContainError(actual interface{}, expected ...interface{}) string {
	if len(expected) != 1 {
		return fmt.Sprintf(
			"expected exactly one argument, got %d",
			len(expected),
		)
	}
	err, ok := actual.(error)
	if !ok {
		return fmt.Sprintf("expected actual to be error, not %T", actual)
	}
	target, ok := expected[0].(error)
	if !ok {
		return fmt.Sprintf("expected argument to be error, not %T", expected[0])
	}
	if errorIs(err, target) {
		return ""
	}
	return fmt.Sprintf("expected %+v to contain %+v", err, target)
}

func errorIs(err, target error) bool {
	var errs *multierror.Error
	if errors.As(err, &errs) {
		for _, err := range errs.Errors {
			if errorIs(err, target) {
				return true
			}
		}
	}
	return errors.Is(err, target)
}
