package rclgo_test

import (
	"context"
	"errors"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy/v2"
	. "github.com/smartystreets/goconvey/convey"
	action_msgs_msg "github.com/tiiuae/rclgo/internal/msgs/action_msgs/msg"
	action_msgs_srv "github.com/tiiuae/rclgo/internal/msgs/action_msgs/srv"
	test_msgs_action "github.com/tiiuae/rclgo/internal/msgs/test_msgs/action"
	"github.com/tiiuae/rclgo/pkg/rclgo"
	"github.com/tiiuae/rclgo/pkg/rclgo/types"
)

var (
	actionServerOpts = rclgo.NewDefaultActionServerOptions()
	actionClientOpts = rclgo.NewDefaultActionClientOptions()
)

func init() {
	actionServerOpts.GoalServiceQos.History = rclgo.RmwQosHistoryPolicyKeepAll
	actionServerOpts.CancelServiceQos.History = rclgo.RmwQosHistoryPolicyKeepAll
	actionServerOpts.ResultServiceQos.History = rclgo.RmwQosHistoryPolicyKeepAll
	actionServerOpts.FeedbackTopicQos = reliableQos
	actionServerOpts.StatusTopicQos = reliableQos

	actionClientOpts.GoalServiceQos.History = rclgo.RmwQosHistoryPolicyKeepAll
	actionClientOpts.CancelServiceQos.History = rclgo.RmwQosHistoryPolicyKeepAll
	actionClientOpts.ResultServiceQos.History = rclgo.RmwQosHistoryPolicyKeepAll
	actionClientOpts.FeedbackTopicQos = reliableQos
	actionClientOpts.StatusTopicQos = reliableQos
}

func newWaitAction() (chan struct{}, rclgo.Action) {
	ch := make(chan struct{})
	return ch, rclgo.NewAction(
		test_msgs_action.FibonacciTypeSupport,
		func(ctx context.Context, goal *rclgo.GoalHandle) (types.Message, error) {
			desc := goal.Description.(*test_msgs_action.Fibonacci_Goal)
			if desc.Order < 0 {
				return nil, errors.New("order is negative")
			}
			if _, err := goal.Accept(); err != nil {
				return nil, err
			}
			if desc.Order > 0 {
				return nil, errors.New("order is positive")
			}
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-ch:
			}
			return test_msgs_action.NewFibonacci_Result(), nil
		},
	)
}

type fibonacciAction struct {
	continueChan chan struct{}
}

func (a *fibonacciAction) TypeSupport() types.ActionTypeSupport {
	return test_msgs_action.FibonacciTypeSupport
}

func (a *fibonacciAction) ExecuteGoal(ctx context.Context, goal *rclgo.GoalHandle) (types.Message, error) {
	desc := goal.Description.(*test_msgs_action.Fibonacci_Goal)
	if desc.Order < 0 {
		return nil, errors.New("order must be non-negative")
	}
	sender, err := goal.Accept()
	if err != nil {
		return nil, err
	}
	if desc.Order >= 100 {
		return nil, errors.New("order must be less than 100")
	}
	result := test_msgs_action.NewFibonacci_Result()
	fb := test_msgs_action.NewFibonacci_Feedback()
	var x, y, i int32
	for y = 1; i < desc.Order; x, y, i = y, x+y, i+1 {
		result.Sequence = append(result.Sequence, x)
		fb.Sequence = result.Sequence
		if err = sender.Send(fb); err != nil {
			goal.Logger().Error("failed to send feedback: ", err)
		} else {
			goal.Logger().Debug("sent feedback: ", i+1)
		}
	}
	sender, err = goal.Accept()
	if err != nil {
		return nil, err
	}
	if sender == nil {
		return nil, errors.New("second call to accept should return a non-nil FeedbackSender")
	}
	<-a.continueChan
	return result, nil
}

type fibonacciFeedbacks []*test_msgs_action.Fibonacci_Feedback

func (a fibonacciFeedbacks) Len() int      { return len(a) }
func (a fibonacciFeedbacks) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a fibonacciFeedbacks) Less(i, j int) bool {
	return len(a[i].Sequence) <= len(a[j].Sequence)
}

func TestActionExecution(t *testing.T) {
	var (
		ctx, cancel  = context.WithCancel(context.Background())
		rclctx       *rclgo.Context
		node1, node2 *rclgo.Node
		client       *rclgo.ActionClient
		err          error
		spinErr      = make(chan error, 2)
		action       = &fibonacciAction{
			continueChan: make(chan struct{}, 10),
		}
	)
	defer func() {
		cancel()
		if rclctx != nil {
			rclctx.Close()
		}
	}()
	Convey("Scenario: ActionClient calls ActionServer and monitors feedback", t, func() {
		Convey("Create an ActionServer", func() {
			rclctx, err = newDefaultRCLContext()
			So(err, ShouldBeNil)
			node1, err = rclctx.NewNode("fibonacci1", "actions_test")
			So(err, ShouldBeNil)
			_, err = node1.NewActionServer("fibonacci", action, actionServerOpts)
			So(err, ShouldBeNil)
		})
		Convey("Create an ActionClient", func() {
			node2, err = rclctx.NewNode("fibonacci2", "actions_test")
			So(err, ShouldBeNil)
			client, err = node2.NewActionClient("fibonacci", test_msgs_action.FibonacciTypeSupport, actionClientOpts)
			So(err, ShouldBeNil)
		})
		Convey("Spin nodes", func() {
			go func() { spinErr <- node1.Spin(ctx) }()
			go func() { spinErr <- node2.Spin(ctx) }()
		})
		Convey("Watch a successful goal", func() {
			goal := test_msgs_action.NewFibonacci_Goal()
			goal.Order = 10
			feedbacks := make(fibonacciFeedbacks, 0)
			var feedbacksMu sync.Mutex
			result, err := client.WatchGoal(ctx, goal, func(c context.Context, m types.Message) {
				fb := m.(*test_msgs_action.Fibonacci_FeedbackMessage)
				feedbacksMu.Lock()
				feedbacks = append(feedbacks, &fb.Feedback)
				if len(feedbacks) == int(goal.Order) {
					action.continueChan <- struct{}{}
				}
				feedbacksMu.Unlock()
			})
			So(err, ShouldBeNil)
			So(result, ShouldNotBeNil)
			sort.Sort(feedbacks)
			So(
				cupaloy.SnapshotMulti(
					"Watch a successful goal",
					result,
					feedbacks,
				),
				ShouldBeNil,
			)
		})
		Convey("Watch a rejected goal", func() {
			close(action.continueChan)
			goal := test_msgs_action.NewFibonacci_Goal()
			goal.Order = -1
			resp, err := client.WatchGoal(ctx, goal, func(c context.Context, m types.Message) {
				panic("no feedback should be sent")
			})
			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})
		Convey("Send a goal that will be rejected", func() {
			goal := test_msgs_action.NewFibonacci_Goal()
			goal.Order = -1
			resp, id, err := client.SendGoal(ctx, goal)
			So(err, ShouldBeNil)
			So(id, ShouldNotBeNil)
			So(resp, ShouldNotBeNil)
			r := resp.(*test_msgs_action.Fibonacci_SendGoal_Response)
			So(r.Accepted, ShouldBeFalse)
		})
		Convey("Send a goal that will be aborted", func() {
			goal := test_msgs_action.NewFibonacci_Goal()
			goal.Order = 100
			resp, id, err := client.SendGoal(ctx, goal)
			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(id, ShouldNotBeNil)
			r := resp.(*test_msgs_action.Fibonacci_SendGoal_Response)
			So(r.Accepted, ShouldBeTrue)
			resp, err = client.GetResult(ctx, id)
			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			result := resp.(*test_msgs_action.Fibonacci_GetResult_Response)
			So(rclgo.GoalStatus(result.Status), ShouldEqual, rclgo.GoalAborted)
		})
		Convey("Resources are released properly", func() {
			cancel()
			timeOut(1000, func() {
				<-spinErr
				<-spinErr
			}, "Waiting for spinning to stop")
			So(rclctx.Close(), ShouldBeNil)
		})
	})
}

func TestActionCanceling(t *testing.T) {
	continueChan, cancelingAction := newWaitAction()
	var (
		ctx, cancel = context.WithCancel(context.Background())
		rclctx      *rclgo.Context
		client      *rclgo.ActionClient
		err         error
		spinErr     = make(chan error, 1)
	)
	defer func() {
		cancel()
		if rclctx != nil {
			rclctx.Close()
		}
	}()
	Convey("Scenario: ActionClient calls ActionServer and monitors feedback", t, func() {
		Convey("Create an ActionServer", func() {
			rclctx, err = newDefaultRCLContext()
			So(err, ShouldBeNil)
			node, err := rclctx.NewNode("cancel1", "actions_test")
			So(err, ShouldBeNil)
			_, err = node.NewActionServer("canceling", cancelingAction, actionServerOpts)
			So(err, ShouldBeNil)
		})
		Convey("Create an ActionClient", func() {
			node, err := rclctx.NewNode("cancel2", "actions_test")
			So(err, ShouldBeNil)
			client, err = node.NewActionClient("canceling", test_msgs_action.FibonacciTypeSupport, actionClientOpts)
			So(err, ShouldBeNil)
		})
		Convey("Spin RCL context", func() {
			go func() { spinErr <- rclctx.Spin(ctx) }()
		})
		Convey("Cancel one goal while another is running", func() {
			goal := test_msgs_action.NewFibonacci_Goal()
			resp, idToCancel, err := client.SendGoal(ctx, goal)
			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(idToCancel, ShouldNotBeNil)
			resp, idToSucceed, err := client.SendGoal(ctx, goal)
			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(idToSucceed, ShouldNotBeNil)
			cancelReq := action_msgs_srv.NewCancelGoal_Request()
			cancelReq.GoalInfo.GoalId.Uuid = *idToCancel
			resp, err = client.CancelGoal(ctx, cancelReq)
			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			resp, err = client.GetResult(ctx, idToCancel)
			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			result := resp.(*test_msgs_action.Fibonacci_GetResult_Response)
			So(rclgo.GoalStatus(result.Status), ShouldEqual, rclgo.GoalCanceled)
			continueChan <- struct{}{}
			resp, err = client.GetResult(ctx, idToSucceed)
			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			result = resp.(*test_msgs_action.Fibonacci_GetResult_Response)
			So(rclgo.GoalStatus(result.Status), ShouldEqual, rclgo.GoalSucceeded)
		})
		Convey("Resources are released properly", func() {
			cancel()
			timeOut(1000, func() {
				err = <-spinErr
			}, "Waiting for spinning to stop")
			So(err, ShouldNotBeNil)
			So(rclctx.Close(), ShouldBeNil)
		})
	})
}

type goalStatus struct {
	ID     types.GoalID
	Status rclgo.GoalStatus
}

type goalStatuses []goalStatus

func (a goalStatuses) Len() int      { return len(a) }
func (a goalStatuses) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a goalStatuses) Less(i, j int) bool {
	return a[i].Status < a[j].Status
}

func TestActionStatuses(t *testing.T) {
	continueChan, statusAction := newWaitAction()
	var (
		ctx, cancel = context.WithCancel(context.Background())
		rclctx      *rclgo.Context
		client      *rclgo.ActionClient
		err         error
		spinErr     = make(chan error, 1)
	)
	defer func() {
		cancel()
		if rclctx != nil {
			rclctx.Close()
		}
	}()
	Convey("Scenario: ActionClient calls ActionServer and monitors feedback", t, func() {
		Convey("Create an ActionServer", func() {
			rclctx, err = newDefaultRCLContext()
			So(err, ShouldBeNil)
			node, err := rclctx.NewNode("status1", "actions_test")
			So(err, ShouldBeNil)
			_, err = node.NewActionServer("statuses", statusAction, actionServerOpts)
			So(err, ShouldBeNil)
		})
		Convey("Create an ActionClient", func() {
			node, err := rclctx.NewNode("status2", "actions_test")
			So(err, ShouldBeNil)
			client, err = node.NewActionClient("statuses", test_msgs_action.FibonacciTypeSupport, actionClientOpts)
			So(err, ShouldBeNil)
		})
		Convey("Spin RCL context", func() {
			go func() { spinErr <- rclctx.Spin(ctx) }()
		})
		//nolint:contextcheck // send returns an inherited context
		Convey("Server reports correct statuses for goals", func() {
			type testResult struct {
				Result   types.Message
				Statuses goalStatuses
				Err      error
			}
			var goalIDCounter uint8
			newGoalID := func() (id types.GoalID) {
				goalIDCounter++
				id[0] = goalIDCounter
				return
			}
			var statuses goalStatuses
			var statusesMu sync.Mutex
			sendCancel := func(ctx context.Context, id *types.GoalID) {
				req := action_msgs_srv.NewCancelGoal_Request()
				req.GoalInfo.GoalId.Uuid = *id
				resp, err := client.CancelGoal(ctx, req)
				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
			}
			send := func(order int32) (*types.GoalID, context.Context, context.CancelFunc) {
				statuses = nil
				ctx, cancel := context.WithCancel(ctx)
				req := test_msgs_action.NewFibonacci_SendGoal_Request()
				req.Goal.Order = order
				req.GoalID.Uuid = newGoalID()
				watchErr := client.WatchStatus(ctx, req.GetGoalID(), func(c context.Context, m types.Message) {
					status := m.(*action_msgs_msg.GoalStatus)
					statusesMu.Lock()
					statuses = append(statuses, goalStatus{
						ID:     *status.GetGoalID(),
						Status: rclgo.GoalStatus(status.Status),
					})
					statusesMu.Unlock()
				})
				So(watchErr, ShouldNotBeNil)
				resp, err := client.SendGoalRequest(ctx, req)
				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
				return req.GetGoalID(), ctx, func() {
					cancel()
					timeOut(1000, func() { err = <-watchErr }, "Wait for watching to stop")
					So(err, ShouldEqual, context.Canceled)
				}
			}
			waitForStatus := func(status rclgo.GoalStatus) {
				timeOut(1000, func() {
					for {
						statusesMu.Lock()
						if len(statuses) > 0 && statuses[len(statuses)-1].Status == status {
							statusesMu.Unlock()
							return
						}
						statusesMu.Unlock()
						time.Sleep(10 * time.Millisecond)
					}
				}, "Wait for status change")
			}
			var testResults []testResult
			addResult := func(ctx context.Context, id *types.GoalID) {
				resp, err := client.GetResult(ctx, id)
				testResults = append(testResults, testResult{
					Result:   resp,
					Statuses: statuses,
					Err:      err,
				})
			}

			id, ctx, cancel := send(0) // should succeed
			continueChan <- struct{}{}
			waitForStatus(rclgo.GoalSucceeded)
			addResult(ctx, id)
			cancel()

			id, ctx, cancel = send(0) // should be canceled
			sendCancel(ctx, id)
			waitForStatus(rclgo.GoalCanceled)
			addResult(ctx, id)
			cancel()

			id, ctx, cancel = send(-1) // should be rejected
			addResult(ctx, id)
			cancel()

			id, ctx, cancel = send(1) // should be aborted
			waitForStatus(rclgo.GoalAborted)
			addResult(ctx, id)
			cancel()

			So(
				cupaloy.SnapshotMulti(
					"Server reports correct statuses for goals",
					testResults,
				),
				ShouldBeNil,
			)
		})
		Convey("Resources are released properly", func() {
			cancel()
			timeOut(1000, func() {
				err = <-spinErr
			}, "Waiting for spinning to stop")
			So(err, ShouldNotBeNil)
			So(rclctx.Close(), ShouldBeNil)
		})
	})
}
