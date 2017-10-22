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
	"fmt"
	"os"
	"runtime"

	"github.com/ckatsak/glocc"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// Print the total results to the standard output in JSON format. This is the
// default and only way to display them for now (other than printing the raw
// Go map, which sucks).
func displayJSON(res glocc.DirResult) {
	if output, err := json.MarshalIndent(res, "", "    "); err != nil {
		fmt.Println(res)
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Println(string(output))
	}
}

func main() {
	totalResults := glocc.DirResult{
		Name:    "TOTAL",
		Subdirs: make(glocc.DirResults, 0),
		Files:   make([]glocc.FileResult, 0),
		Summary: make(map[string]int),
	}
	resultsChannel := make(chan glocc.DirResult)
	for _, path := range os.Args[1:] {
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
		if resultsCount == len(os.Args)-1 {
			break
		}
	}
	displayJSON(totalResults)
}
