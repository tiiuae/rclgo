package rclgo_test

func onErr(err *error, f func() error) {
	if *err != nil {
		f()
	}
}
