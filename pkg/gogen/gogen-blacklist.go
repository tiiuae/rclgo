/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package gogen

import (
	"strings"

	"github.com/kivilahtio/go-re/v0"
)

var blacklistMatchingRegexp string

func init() {
	prepareBlacklistMatchingRegexp()
}
func prepareBlacklistMatchingRegexp() {
	blacklistMatchingRegexp = "(" + strings.Join(ROS2_MESSAGES_BLACKLIST, ")|(") + ")"
	re.S(&blacklistMatchingRegexp, `s!\!!\!!`)
	blacklistMatchingRegexp = "m!" + blacklistMatchingRegexp + "!"
}

func blacklisted(path string) (bool, string) {
	if r := re.Mr(path, blacklistMatchingRegexp); r.Matches > 0 {
		var matchedBlackListRegex string
		for i := 1; i < len(r.S); i++ {
			if r.S[i] != "" {
				matchedBlackListRegex = ROS2_MESSAGES_BLACKLIST[i-1]
			}
		}
		return true, matchedBlackListRegex
	}
	return false, ""
}
