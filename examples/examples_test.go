package test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func skipIfShort(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
}

func shell(t *testing.T, script string) {
	t.Helper()
	cmd := exec.Command("bash", "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	require.NoError(t, cmd.Run())
}

func TestCustomMessagePackage(t *testing.T) {
	skipIfShort(t)
	shell(t, `
set -e
function kill-jobs {
	kill $(jobs -p)
}
trap 'kill-jobs' SIGINT SIGTERM EXIT

cd custom_message_package/greeting_msgs
rm -rf build install log
colcon build
. install/local_setup.sh

cd ../greeter
rm -rf msgs greeter
go generate
go build

./greeter &
msg=$(ros2 topic echo --once /greeter/hello greeting_msgs/msg/Hello)
test "$msg" == $'name: gopher\n---'
`)
}
