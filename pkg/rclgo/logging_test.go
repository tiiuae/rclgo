package rclgo

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInitLogging(t *testing.T) {
	var (
		defaultLogger, testLogger *Logger
		parentLogger, childLogger *Logger
	)
	Convey("Scenario: Logging cannot be used before calling InitLogging", t, func() {
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
		Convey("Loggers don't return an error after initialization", func() {
			args, _, err := ParseArgs([]string{"--ros-args", "--log-level", "DEBUG"})
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
		})
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
			ctx, err = NewContext(nil, 0, args)
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
