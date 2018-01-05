/*
Copyright (C) 2017, Christos Katsakioris
All rights reserved.

This software may be modified and distributed under the terms
of the BSD 3-Clause License. See the LICENSE file for details.
*/

// The glocc command line tool.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/ckatsak/glocc"
	"golang.org/x/sys/unix"
	"gopkg.in/yaml.v2"
)

// Command line flags.
var (
	debugFlag, showAllFlag, showTimeFlag *bool
	outFormatFlag                        *string
)

// Set the soft limit of RLIMIT_NOFILE to be equal to the hard limit, to allow
// as many open files as possible. (How many? Check `/proc/<PID>/limits` to see
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

// Print the total results to the standard output in raw Go map %#v format.
func displayRaw(res interface{}) {
	fmt.Printf("%#v\n", res)
}

// Print the total results to the standard output in JSON format. It falls back
// to printing the raw Go map in case of a failure.
func displayJSON(res interface{}) {
	if output, err := json.MarshalIndent(res, "", "   "); err != nil {
		displayRaw(res)
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Println(string(output))
	}
}

// Print the total results to the standard output in YAML format. It falls back
// to displayJSON in case of failure during marshalling.
func displayYAML(res interface{}) {
	if output, err := yaml.Marshal(res); err != nil {
		displayJSON(res)
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Println(string(output))
	}
}

// It receives a slice of strings, the command line arguments of glocc, and
// returns the total results of counting using the glocc package.
func gloccMain(args []string) glocc.DirResult {
	totalResults := glocc.DirResult{
		Name:    "TOTAL",
		Subdirs: make(glocc.DirResults, 0),
		Files:   make([]glocc.FileResult, 0),
		Summary: make(map[string]int),
	}
	resultsChannel := make(chan glocc.DirResult)
	for _, path := range args {
		go func(path string) {
			resultsChannel <- glocc.CountLoc(path)
		}(path)
	}
	resultsCount := 0
	for result := range resultsChannel {
		resultsCount++
		totalResults.Subdirs = append(totalResults.Subdirs, result)
		for lang, loc := range result.Summary {
			if _, exists := totalResults.Summary[lang]; exists {
				totalResults.Summary[lang] += loc
			} else {
				totalResults.Summary[lang] = loc
			}
		}
		if resultsCount == len(args) {
			break
		}
	}
	return totalResults
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	debugFlag = flag.Bool("debug", false, "enable verbose logging to standard error; useful for debugging")
	showAllFlag = flag.Bool("a", false, "show extensive results instead of just a top-level summary (default is summary)")
	outFormatFlag = flag.String("o", "yaml", "choose output format; YAML, JSON and \"raw\" are currently supported")
	showTimeFlag = flag.Bool("t", false, "print the total duration of counting all arguments")
}

func main() {
	flag.Parse()

	if *debugFlag {
		glocc.EnableLogging()
	}

	var displayFunc func(interface{})
	switch strings.ToLower(*outFormatFlag) {
	case "json":
		displayFunc = displayJSON
	case "yaml", "yml":
		displayFunc = displayYAML
	case "raw":
		displayFunc = displayRaw
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	setNoFilesHardLimit()

	startTime := time.Now()
	totalResults := gloccMain(flag.Args())
	endTime := time.Since(startTime)

	if *showAllFlag {
		displayFunc(totalResults)
	} else {
		displayFunc(totalResults.Summary)
	}

	if *showTimeFlag {
		fmt.Printf("Counting completed in %s.\n", endTime)
	}
}
