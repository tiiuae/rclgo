package ros2

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	std_msgs "github.com/tiiuae/rclgo/pkg/ros2/msgs/std_msgs/msg"
	test_msgs "github.com/tiiuae/rclgo/pkg/ros2/msgs/test_msgs/msg"
)

func TestPubSub(t *testing.T) {
	var rclContextPub *Context
	var rclContextSub *Context
	var waitSet *WaitSet
	var publisher *Publisher
	var errsSub, errsPub error
	subChan := make(chan *std_msgs.ColorRGBA, 1)
	subCtx, cancelSubCtx := context.WithCancel(context.Background())
	defer cancelSubCtx()
	SetDefaultFailureMode(FailureContinues)
	Convey("Scenario: Publisher publishes, Subscriber subscribes, garbage is collected", t, func() {

		Convey("Given a Subscriber", func() {
			rclContextSub, waitSet, errsSub = SubscriberBundleReturnWaitSet(subCtx, rclContextSub, nil, "/test", "", "/topic", "std_msgs/ColorRGBA", NewRCLArgsMust("--ros-args --log-level DEBUG"),
				func(s *Subscription) {
					var m std_msgs.ColorRGBA
					if _, err := s.TakeMessage(&m); err != nil {
						fmt.Println("failed to take message:", err)
					}
					subChan <- &m
				})
			So(errsSub, ShouldBeNil)
		})
		Convey("And a Publisher", func() {
			rclContextPub, publisher, errsPub = PublisherBundle(rclContextPub, nil, "/test", "", "/topic", "std_msgs/ColorRGBA", NewRCLArgsMust("--ros-args --log-level DEBUG"))
			So(errsPub, ShouldBeNil)
		})
		Convey("And the Subscriber is ready to work", func() {
			err := waitSet.WaitForReady(5*time.Second, 10*time.Millisecond)
			So(err, ShouldBeNil)
		})
		Convey("When the Publisher publishes", func() {
			err := publishColorRGBA(publisher, 1.55, 2.66, 3.77, 4.88)
			So(err, ShouldBeNil)
		})
		Convey("Then the Subscriber receives", func() {
			receiveColorRGBA(rclContextSub, subChan, 1.55, 2.66, 3.77, 4.88)
		})
		Convey("When the Publisher publishes again", func() {
			err := publishColorRGBA(publisher, 0.00, 1.00, 2.00, 3.00)
			So(err, ShouldBeNil)
		})
		Convey("Then the Subscriber receives again", func() {
			receiveColorRGBA(rclContextSub, subChan, 0.00, 1.00, 2.00, 3.00)
		})
		Convey("When the Subscriber context is canceled", func() {
			cancelSubCtx()
			timeOut(2000, func() { rclContextSub.WG.Wait() }, "Subscriber waitGroup waiting to finish")
		})
		Convey("And the Subscriber context is GC'd", func() {
			errs := rclContextSub.Close()
			So(errs, ShouldBeNil)
		})
		Convey("And the Publisher publishes to a Topic with no Subscribers", func() {
			err := publishColorRGBA(publisher, 0.00, 1.00, 2.00, 3.00)
			So(err, ShouldBeNil)
		})
		Convey("Then the Subscriber cannot receive the message", func() {
			So(len(subChan), ShouldEqual, 0)
		})
		Convey("And the Publisher context is GC'd", func() {
			errs := rclContextPub.Close()
			So(errs, ShouldBeNil)
		})
	})
}

func TestMultipleSubscribersInSingleWaitSet(t *testing.T) {
	var (
		rclCtxPub, rclCtxSub *Context
		pub1, pub2           *Publisher
		sub1, sub2           *Subscription
		topic1Chan           = make(chan receiveResult, 1)
		topic2Chan           = make(chan receiveResult, 1)
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
	Convey("Scenario: Publishers publishe, multiple Subscribers listen", t, func() {
		Convey("Given RCL contexts and waitset", func() {
			var err error
			rclCtxSub, err = newDefaultRCLContext()
			So(err, ShouldBeNil)
			rclCtxPub, err = newDefaultRCLContext()
			So(err, ShouldBeNil)
		})
		Convey("And a subscriber on the first topic", func() {
			node, err := rclCtxPub.NewNode("sub1", "/test")
			So(err, ShouldBeNil)
			sub1, err = node.NewSubscription("/topic1", std_msgs.StringTypeSupport, sendToChan(topic1Chan))
			So(err, ShouldBeNil)
		})
		Convey("And a subscriber on the second topic", func() {
			node, err := rclCtxPub.NewNode("sub2", "/test")
			So(err, ShouldBeNil)
			sub2, err = node.NewSubscription("/topic2", std_msgs.StringTypeSupport, sendToChan(topic2Chan))
			So(err, ShouldBeNil)
		})
		Convey("And a publisher on the first topic", func() {
			node, err := rclCtxPub.NewNode("pub1", "/test")
			So(err, ShouldBeNil)
			pub1, err = node.NewPublisher("/topic1", std_msgs.StringTypeSupport)
			So(err, ShouldBeNil)
		})
		Convey("And a publisher on the second topic", func() {
			node, err := rclCtxPub.NewNode("pub2", "/test")
			So(err, ShouldBeNil)
			pub2, err = node.NewPublisher("/topic1", std_msgs.StringTypeSupport)
			So(err, ShouldBeNil)
		})
		Convey("And the waitset is started", func() {
			subWaitSet, err := rclCtxSub.NewWaitSet(time.Second)
			So(err, ShouldBeNil)
			subWaitSet.AddSubscriptions(sub1, sub2)
			subWaitSet.RunGoroutine(ctx)
			So(subWaitSet.WaitForReady(5*time.Second, 10*time.Millisecond), ShouldBeNil)
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
			timeOut(2000, func() { rclCtxSub.WG.Wait() }, "Subscriber waitGroup waiting to finish")
		})
		Convey("And RCL contexts are freed", func() {
			So(rclCtxSub.Close(), ShouldBeNil)
			So(rclCtxPub.Close(), ShouldBeNil)
		})
	})
}

func TestMultipleTimersInSingleWaitSet(t *testing.T) {
	var (
		rclCtx         *Context
		timer1, timer2 *Timer
		timer1Chan     = make(chan struct{}, 1)
		timer2Chan     = make(chan struct{}, 1)
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
			timer1, err = rclCtx.NewTimer(time.Second, func(t *Timer) {
				timer1Chan <- struct{}{}
			})
			So(err, ShouldBeNil)
			timer2, err = rclCtx.NewTimer(time.Hour, func(t *Timer) {
				timer2Chan <- struct{}{}
			})
			So(err, ShouldBeNil)
		})
		Convey("When WaitSet is started", func() {
			waitSet, err := rclCtx.NewWaitSet(500 * time.Millisecond)
			So(err, ShouldBeNil)
			waitSet.AddTimers(timer1, timer2)
			waitSet.RunGoroutine(ctx)
			So(waitSet.WaitForReady(5*time.Second, 10*time.Millisecond), ShouldBeNil)
		})
		Convey("Then callback of timer 1 is called", func() {
			timeOut(1050, func() { <-timer1Chan }, "Waiting for timer 1 callback")
		})
		Convey("But callback of timer 2 is not called", func() {
			receiveNothing(timer2Chan)
		})
		Convey("When the context is canceled", func() {
			cancel()
		})
		Convey("Then all subscribers stop", func() {
			timeOut(2000, func() { rclCtx.WG.Wait() }, "Subscriber waitGroup waiting to finish")
		})
		Convey("And RCL contexts are freed", func() {
			So(rclCtx.Close(), ShouldBeNil)
		})
	})
}

func BenchmarkPubsubMemoryLeakAllocateInLoop(t *testing.B) {
	var messagesReceived int = 0
	fmt.Printf("Mem from pmap(1) '%skB' messages '%d'\n", getMemReading(), messagesReceived)
	for {
		runCtx, stopRun := context.WithCancel(context.Background())
		rclContextSub, waitSet, errs := SubscriberBundleReturnWaitSet(runCtx, nil, nil, "/test", "", "/topic", "test_msgs/UnboundedSequences", nil,
			func(s *Subscription) {
				var m test_msgs.UnboundedSequences
				_, err := s.TakeMessage(&m)
				if err != nil {
					fmt.Println("failed to take message:", err)
				}
				//fmt.Printf("%+v\n", c)
				messagesReceived++
			})
		if errs != nil {
			fmt.Println("error:", errs)
			continue
		}
		err := waitSet.WaitForReady(2*time.Second, 10*time.Millisecond)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}

		rclContextPub, errs := PublisherBundleTimer(runCtx, nil, nil, "/test", "", "/topic", "test_msgs/UnboundedSequences", nil, 1*time.Millisecond, "", nil)
		if errs != nil {
			fmt.Println("error:", errs)
			continue
		}

		time.Sleep(1000 * time.Millisecond)

		stopRun()
		errs = rclContextSub.Close()
		if errs != nil {
			fmt.Println("error:", errs)
		}
		errs = rclContextPub.Close()
		if errs != nil {
			fmt.Println("error:", errs)
		}
		runtime.GC()
		fmt.Printf("Mem from pmap(1) '%skB' messages '%d'\n", getMemReading(), messagesReceived)
	}
}

func BenchmarkPubsubMemoryLeakAllocateOutOfLoop(t *testing.B) {
	var messagesReceived int = 0
	fmt.Printf("Mem from pmap(1) '%skB' messages '%d'\n", getMemReading(), messagesReceived)
	runCtx, stopRun := context.WithCancel(context.Background())
	rclContext, waitSet, errs := SubscriberBundleReturnWaitSet(runCtx, nil, nil, "/test", "", "/topic", "test_msgs/UnboundedSequences", nil,
		func(s *Subscription) {
			var m test_msgs.UnboundedSequences
			_, err := s.TakeMessage(&m)
			if err != nil {
				fmt.Println("failed to take message:", err)
			}
			//fmt.Printf("%+v\n", c)
			messagesReceived++
		})
	if errs != nil {
		panic(errs)
	}

	err := waitSet.WaitForReady(2*time.Second, 10*time.Millisecond)
	if err != nil {
		panic(err)
	}

	rclContext, errs = PublisherBundleTimer(runCtx, rclContext, nil, "/test", "", "/topic", "test_msgs/UnboundedSequences", nil, 1*time.Millisecond, "", nil)
	if errs != nil {
		panic(errs)
	}
	defer rclContext.Close()
	defer stopRun()

	for {
		time.Sleep(1000 * time.Millisecond)
		runtime.GC()
		fmt.Printf("Mem from pmap(1) '%skB' messages '%d'\n", getMemReading(), messagesReceived)
	}
}

func getMemReading() string {
	cmd := `pmap ` + fmt.Sprint(os.Getpid()) + ` | tail -n 1 | grep -Po '\d+'` //  total          2102728K => 2102728
	output, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return fmt.Sprintf("Failed to execute command: %s", cmd)
	}
	return strings.TrimSpace(string(output))
}

func publishColorRGBA(p *Publisher, r, g, b, a float32) error {
	m := p.typeSupport.New().(*std_msgs.ColorRGBA)
	m.R = r
	m.G = g
	m.B = b
	m.A = a
	return p.Publish(m)
}

func receiveColorRGBA(c *Context, subChan chan *std_msgs.ColorRGBA, r, g, b, a float32) {
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

func publishString(pub *Publisher, s string) {
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

func newDefaultRCLContext() (*Context, error) {
	return NewContext(&sync.WaitGroup{}, 0, defaultRCLArgs())
}

func defaultRCLArgs() *RCLArgs {
	osArgs := os.Args
	defer func() { os.Args = osArgs }()
	os.Args = []string{}
	return NewRCLArgsMust("")
}

func sendToChan(c chan<- receiveResult) func(s *Subscription) {
	return func(s *Subscription) {
		var res receiveResult
		res.rmi, res.err = s.TakeMessage(&res.msg)
		c <- res
	}
}

type receiveResult struct {
	msg std_msgs.String
	rmi *RmwMessageInfo
	err error
}
