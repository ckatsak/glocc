// Copyright 2018 Christos Katsakioris
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
