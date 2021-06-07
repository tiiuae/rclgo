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

/*
prepareBlacklistMatchingRegexp is a convenience function to be more easily able to define the files to blacklist without needing to fiddle with complex regexp
*/
func prepareBlacklistMatchingRegexp() {
	blacklistMatchingRegexp = "(" + strings.Join(blacklistedMessages, ")|(") + ")"
	re.S(&blacklistMatchingRegexp, `s!\!!\!!`)
	blacklistMatchingRegexp = "m!" + blacklistMatchingRegexp + "!"
}

func blacklisted(path string) (bool, string) {
	if r := re.Mr(path, blacklistMatchingRegexp); r.Matches > 0 {
		var matchedBlackListRegex string
		for i := 1; i < len(r.S); i++ {
			if r.S[i] != "" {
				matchedBlackListRegex = blacklistedMessages[i-1]
			}
		}
		return true, matchedBlackListRegex
	}
	return false, ""
}
