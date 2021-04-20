package test

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/tiiuae/rclgo/pkg/ros2"
	std_msgs "github.com/tiiuae/rclgo/pkg/ros2/msgs/std_msgs/msg"
	test_msgs "github.com/tiiuae/rclgo/pkg/ros2/msgs/test_msgs/msg"
)

func TestPubSub(t *testing.T) {
	var rclContextPub ros2.RCLContext
	var rclContextSub ros2.RCLContext
	var errsSub, errsPub *ros2.RCLErrors
	subChan := make(chan *std_msgs.ColorRGBA, 1)
	SetDefaultFailureMode(FailureContinues)
	Convey("Scenario: Publisher publishes, Subscriber subscribes, garbage is collected", t, func() {

		Convey("Given a Subscriber", func() {
			rclContextSub, errsSub = ros2.SubscriberBundle(rclContextSub, nil, "/test", "", "/topic", "std_msgs/ColorRGBA", ros2.NewRCLArgsMust("--ros-args --log-level DEBUG"),
				func(s *ros2.Subscription, p unsafe.Pointer, rmi *ros2.RmwMessageInfo) {
					m := s.Ros2MsgType.Clone()
					m.AsGoStruct(p)
					c := m.(*std_msgs.ColorRGBA)
					subChan <- c
				})
			So(errsSub, ShouldBeNil)
		})
		Convey("And a Publisher", func() {
			rclContextPub, _, errsPub = ros2.PublisherBundle(rclContextPub, nil, "/test", "", "/topic", "std_msgs/ColorRGBA", ros2.NewRCLArgsMust("--ros-args --log-level DEBUG"))
			So(errsPub, ShouldBeNil)
		})
		Convey("And the Subscriber is ready to work", func() {
			w := ros2.GetRCLEntities(rclContextSub).WaitSets.Front().Value.(*ros2.WaitSet)
			err := w.WaitForReady(5000, 10)
			So(err, ShouldBeNil)
		})
		Convey("When the Publisher publishes", func() {
			err := publishColorRGBA(rclContextPub, 1.55, 2.66, 3.77, 4.88)
			So(err, ShouldBeNil)
		})
		Convey("Then the Subscriber receives", func() {
			receiveColorRGBA(rclContextSub, subChan, 1.55, 2.66, 3.77, 4.88)
		})
		Convey("When the Publisher publishes again", func() {
			err := publishColorRGBA(rclContextPub, 0.00, 1.00, 2.00, 3.00)
			So(err, ShouldBeNil)
		})
		Convey("Then the Subscriber receives again", func() {
			receiveColorRGBA(rclContextSub, subChan, 0.00, 1.00, 2.00, 3.00)
		})
		Convey("When the Subscriber context is canceled", func() {
			ros2.GetRCLContextImpl(rclContextSub).Quit()
			timeOut(1000, func() { ros2.GetRCLContextImpl(rclContextSub).WG.Wait() }, "Subscriber waitGroup waiting to finish")
		})
		Convey("And the Subscriber context is GC'd", func() {
			errs := ros2.RCLContextFini(rclContextSub)
			So(errs, ShouldBeNil)
		})
		Convey("And the Publisher publishes to a Topic with no Subscribers", func() {
			err := publishColorRGBA(rclContextPub, 0.00, 1.00, 2.00, 3.00)
			So(err, ShouldBeNil)
		})
		Convey("Then the Subscriber cannot receive the message", func() {
			So(len(subChan), ShouldEqual, 0)
		})
		Convey("Finally the Publisher context is canceled", func() {
			ros2.GetRCLContextImpl(rclContextPub).Quit()
		})
		Convey("And the Publisher context is GC'd", func() {
			errs := ros2.RCLContextFini(rclContextPub)
			So(errs, ShouldBeNil)
		})
	})
}

func BenchsittingmarkMemoryLeak(t *testing.B) {
	var messagesReceived int = 0
	fmt.Printf("Mem from pmap(1) '%skB' messages '%d'\n", getMemReading(), messagesReceived)
	for {
		rclContextSub, errs := ros2.SubscriberBundle(nil, nil, "/test", "", "/topic", "test_msgs/UnboundedSequences", nil,
			func(s *ros2.Subscription, p unsafe.Pointer, rmi *ros2.RmwMessageInfo) {
				m := s.Ros2MsgType.Clone()
				m.AsGoStruct(p)
				_ = m.(*test_msgs.UnboundedSequences)
				//fmt.Printf("%+v\n", c)
				messagesReceived++
			})
		if errs != nil {
			panic(errs)
		}

		err := ros2.GetRCLEntities(rclContextSub).WaitSets.Front().Value.(*ros2.WaitSet).WaitForReady(1000, 10)
		if err != nil {
			panic(err)
		}

		rclContextPub, errs := ros2.PublisherBundleTimer(nil, nil, "/test", "", "/topic", "test_msgs/UnboundedSequences", nil, 1*time.Millisecond, "", nil)
		if errs != nil {
			panic(errs)
		}

		time.Sleep(1000 * time.Millisecond)

		errs = ros2.RCLContextFini(rclContextSub)
		if errs != nil {
			panic(errs)
		}
		errs = ros2.RCLContextFini(rclContextPub)
		if errs != nil {
			panic(errs)
		}

		fmt.Printf("Mem from pmap(1) '%skB' messages '%d'\n", getMemReading(), messagesReceived)
	}
}

func BenchmarkMemoryLeak(t *testing.B) {
	var messagesReceived int = 0
	fmt.Printf("Mem from pmap(1) '%skB' messages '%d'\n", getMemReading(), messagesReceived)
	rclContext, errs := ros2.SubscriberBundle(nil, nil, "/test", "", "/topic", "test_msgs/UnboundedSequences", nil,
		func(s *ros2.Subscription, p unsafe.Pointer, rmi *ros2.RmwMessageInfo) {
			m := s.Ros2MsgType.Clone()
			m.AsGoStruct(p)
			_ = m.(*test_msgs.UnboundedSequences)
			//fmt.Printf("%+v\n", c)
			messagesReceived++
		})
	if errs != nil {
		panic(errs)
	}

	err := ros2.GetRCLEntities(rclContext).WaitSets.Front().Value.(*ros2.WaitSet).WaitForReady(1000, 10)
	if err != nil {
		panic(err)
	}

	rclContext, errs = ros2.PublisherBundleTimer(rclContext, nil, "/test", "", "/topic", "test_msgs/UnboundedSequences", nil, 1*time.Millisecond, "", nil)
	if errs != nil {
		panic(errs)
	}

	for {
		time.Sleep(1000 * time.Millisecond)
		runtime.GC()
		fmt.Printf("Mem from pmap(1) '%skB' messages '%d'\n", getMemReading(), messagesReceived)
	}

	errs = ros2.RCLContextFini(rclContext)
	if errs != nil {
		panic(errs)
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

func publishColorRGBA(c ros2.RCLContext, r, g, b, a float32) ros2.RCLError {
	p := ros2.GetRCLEntities(c).Publishers.Front().Value.(*ros2.Publisher)
	m := p.Ros2MsgType.Clone().(*std_msgs.ColorRGBA)
	m.R = r
	m.G = g
	m.B = b
	m.A = a
	return p.Publish(m)
}

func receiveColorRGBA(c ros2.RCLContext, subChan chan *std_msgs.ColorRGBA, r, g, b, a float32) {
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
