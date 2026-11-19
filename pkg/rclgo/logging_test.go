package rclgo

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestInitLogging(t *testing.T) {
	var (
		defaultLogger, testLogger *Logger
		parentLogger, childLogger *Logger
	)
	Convey(
		"Scenario: Logging cannot be used before calling InitLogging",
		t,
		func() {
			Convey("Loggers can be created before initialization", func() {
				defaultLogger = GetLogger("")
				So(defaultLogger, ShouldNotBeNil)
				testLogger = GetLogger("test.logger")
				So(testLogger, ShouldNotBeNil)
				parentLogger = testLogger.Parent()
				So(parentLogger, ShouldNotBeNil)
				So(parentLogger.Name(), ShouldEqual, "test")
				childLogger = testLogger.Child("child")
				So(childLogger, ShouldNotBeNil)
				So(childLogger.Name(), ShouldEqual, "test.logger.child")
				defaultLogger2 := parentLogger.Parent()
				So(defaultLogger, ShouldEqual, defaultLogger2)
				defaultLogger2 = parentLogger.Parent()
				So(defaultLogger, ShouldEqual, defaultLogger2)
			})
			Convey(
				"Loggers don't return an error after initialization",
				func() {
					args, _, err := ParseArgs(
						[]string{"--ros-args", "--log-level", "DEBUG"},
					)
					So(args, ShouldNotBeNil)
					So(err, ShouldBeNil)

					level, err := childLogger.EffectiveLevel()
					So(level, ShouldEqual, LogSeverityInfo)
					So(err, ShouldBeNil)
					level, err = childLogger.Level()
					So(level, ShouldEqual, LogSeverityUnset)
					So(err, ShouldBeNil)

					So(InitLogging(args), ShouldBeNil)

					level, err = childLogger.EffectiveLevel()
					So(level, ShouldEqual, LogSeverityDebug)
					So(err, ShouldBeNil)
					level, err = childLogger.Level()
					So(level, ShouldEqual, LogSeverityUnset)
					So(err, ShouldBeNil)

					So(childLogger.SetLevel(LogSeverityError), ShouldBeNil)
					level, err = childLogger.Level()
					So(level, ShouldEqual, LogSeverityError)
					So(err, ShouldBeNil)

					So(defaultLogger.Info("defaultLogger"), ShouldBeNil)
					So(testLogger.Info("testLogger"), ShouldBeNil)
					So(parentLogger.Info("parentLogger"), ShouldBeNil)
					So(childLogger.Info("childLogger"), ShouldBeNil)
				},
			)
		},
	)
}

func TestGetLogger(t *testing.T) {
	Convey("Scenario: GetLogger returns correct loggers", t, func() {
		So(GetLogger(""), ShouldEqual, defaultLogger)
		So(GetLogger(""), ShouldEqual, GetLogger(""))

		a := GetLogger("a")
		So(a, ShouldNotBeNil)
		So(a.Name(), ShouldEqual, GetLogger("a").Name())
		So(a.Parent(), ShouldEqual, GetLogger(""))

		ab := GetLogger("a.b")
		So(ab, ShouldNotBeNil)
		So(ab.Name(), ShouldEqual, GetLogger("a.b").Name())
		So(ab.Name(), ShouldNotEqual, a.Name())
		So(ab.Parent().Name(), ShouldEqual, a.Name())
		So(a.Child("b").Name(), ShouldEqual, ab.Name())
		So(a.Child("."), ShouldBeNil)
		So(a.Child(".b"), ShouldBeNil)
		So(a.Child("b."), ShouldBeNil)

		So(GetLogger("."), ShouldBeNil)
		So(GetLogger(".."), ShouldBeNil)
		So(GetLogger("..."), ShouldBeNil)
		So(GetLogger(".a."), ShouldBeNil)
		So(GetLogger(".a.b"), ShouldBeNil)
		So(GetLogger("a.b."), ShouldBeNil)
		So(GetLogger("a..b"), ShouldBeNil)
	})
}

func TestNodeLogger(t *testing.T) {
	var (
		ctx    *Context
		logger *Logger
	)
	defer func() {
		if ctx != nil {
			ctx.Close()
		}
	}()
	Convey("Scenario: Node logger works", t, func() {
		Convey("NewContext initializes logging", func() {
			logger = GetLogger("")
			So(logger, ShouldNotBeNil)

			args, _, err := ParseArgs(nil)
			So(args, ShouldNotBeNil)
			So(err, ShouldBeNil)
			ctx, err = NewContext(0, args)
			So(ctx, ShouldNotBeNil)
			So(err, ShouldBeNil)

			So(logger.Debugln("after NewContext"), ShouldBeNil)
		})
		Convey("Node logger works", func() {
			node, err := ctx.NewNode("logging_node", "test_namespace")
			So(node, ShouldNotBeNil)
			So(err, ShouldBeNil)

			So(node.Logger(), ShouldNotBeNil)
			So(node.Logger().Name(), ShouldEqual, "test_namespace.logging_node")
			So(node.Logger().Info("node logger"), ShouldBeNil)
		})
	})
}

// Reconfiguring at DEBUG makes rcl log from inside rcl_logging_configure through
// rclgo's handler while loggingMutex is held, which used to self-deadlock.
func TestReinitLoggingAtDebugDoesNotDeadlock(t *testing.T) {
	var parsed [3]*Args
	for i, a := range [][]string{
		nil,
		{"--ros-args", "--log-level", "DEBUG"},
		{"--ros-args", "--log-level", "INFO"},
	} {
		args, _, err := ParseArgs(a)
		if err != nil {
			t.Fatal(err)
		}
		parsed[i] = args
	}
	defaults, debug, info := parsed[0], parsed[1], parsed[2]

	done := make(chan error, 1)
	go func() {
		if err := InitLogging(defaults); err != nil {
			done <- err
			return
		}
		done <- InitLogging(debug)
	}()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("InitLogging deadlocked when reconfiguring at DEBUG")
	}
	if err := InitLogging(info); err != nil {
		t.Fatal(err)
	}
}
