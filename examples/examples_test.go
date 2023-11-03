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

const shellPrelude = `
set -eo pipefail
function kill-jobs {
	kill $(jobs -p)
}
trap 'kill-jobs' SIGINT SIGTERM EXIT
`

func shell(t *testing.T, script string) {
	t.Helper()
	cmd := exec.Command("bash", "-c", shellPrelude+script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	require.NoError(t, cmd.Run())
}

func TestCustomMessagePackage(t *testing.T) {
	skipIfShort(t)
	shell(t, `
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

func TestPublisherSubscriber(t *testing.T) {
	skipIfShort(t)
	shell(t, `
cd publisher_subscriber
rm -rf msgs pub sub pipe
go generate

go build -o pub ./publisher
go build -o sub ./subscriber

./pub &

mkfifo pipe
./sub >pipe 2>&1 &

function check_output {
	received_msg='Received: &std_msgs_msg.String{Data:"gopher"}'
	while read -r line; do
		if [[ "$line" == *"$received_msg"* ]]; then
			exit 0
		fi
	done <pipe
	exit 1
}
export -f check_output

timeout 3 bash -c 'check_output'
`)
}
