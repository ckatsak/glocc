// Copyright (C) 2017, Christos Katsakioris
// All rights reserved.
//
// This software may be modified and distributed under the terms
// of the BSD 3-Clause License. See the LICENSE file for details.

// +build !windows

package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// Set the soft limit of RLIMIT_NOFILE to be equal to the hard limit, to allow
// as many open files as possible. (How many? Check /proc/<PID>/limits to see
// for yourself.)
func setNoFilesHardLimit() {
	var rlimit unix.Rlimit
	if err := unix.Getrlimit(unix.RLIMIT_NOFILE, &rlimit); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	rlimit.Cur = rlimit.Max
	if err := unix.Setrlimit(unix.RLIMIT_NOFILE, &rlimit); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
