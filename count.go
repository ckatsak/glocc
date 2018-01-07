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

package glocc

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// DirResult is a tree-like (thus recursive) data structure used to store the
// results of the count for all files and subdirectories that live under the
// directory associated with a directory.
//
// A DirResult contains the following fields:
//
// - Name is the full name of the subdirectory it represents, as a string.
//
// - Subdirs is a slice of DirResult. Each element in the slice, represents the
// results of counting lines of code in a subdirectory under the directory
// associated with this DirResult.
//
// - Files is a slice of FileResult. Each element in the slice represents the
// results of counting lines of code in a file living under the directory
// associated with this DirResult.
//
// - Summary provides a summary of the results of the counting.
type DirResult struct {
	Name    string         `json:"name" yaml:"Name"`
	Subdirs DirResults     `json:"subdirs,omitempty" yaml:"subdirs,omitempty"`
	Files   []FileResult   `json:"files,omitempty" yaml:"files,omitempty"`
	Summary map[string]int `json:"summary" yaml:"Summary"`
}

// DirResults is a slice of DirResult.
type DirResults []DirResult

// FileResult is a simple data structure used to store the results of a single
// file's count. FileResult structs typically live inside DirResult structs.
type FileResult struct {
	Name string         `json:"name" yaml:"Name,omitempty"`
	Loc  map[string]int `json:"loc" yaml:"loc,omitempty,inline"`
}

// Package-level logger.
var logger *log.Logger

func init() {
	logger = log.New(ioutil.Discard, "glocc: ", log.Ltime|log.Lmicroseconds|log.Lshortfile)
}

// EnableLogging enables verbose logging to standard error stream using a
// package-level logger.
// This might be useful for debugging.
func EnableLogging() {
	logger.SetOutput(os.Stderr)
}

// DisableLogging disables verbose logging to standard error stream using the
// package-level logger.
func DisableLogging() {
	logger.SetOutput(ioutil.Discard)
}

// CountLoc is the main exported interface of glocc package, meant to be called
// once for each top-level directory in which counting lines of code is needed.
// It returns a DirResult that contains the results of the counting.
func CountLoc(root string) DirResult {
	start := time.Now()
	result := DirResult{
		Name:    root,
		Subdirs: make(DirResults, 0),
		Files:   make([]FileResult, 0),
		Summary: make(map[string]int),
	}
	rootPath, err := filepath.Abs(root)
	if err != nil {
		logger.Println("ERROR", err)
		return result
	}
	fileinfo, err := os.Stat(rootPath)
	if err != nil {
		logger.Println("ERROR", err)
		return result
	}
	if fileinfo.IsDir() {
		result = locDir(rootPath)
	} else if fileinfo.Mode().IsRegular() {
		fileResult := locFile(rootPath)
		result.Name = fileResult.Name
		result.Subdirs = nil
		result.Files = []FileResult{*fileResult}
		result.Summary = fileResult.Loc
	}
	logger.Printf("INFO Time elapsed for %q: %s\n", root, time.Since(start))
	return result
}

// The core recursive function for diving into subdirectories, and for spawning
// (per file and per subdirectory) and synchronizing the goroutines.
func locDir(rootPath string) DirResult {
	result := DirResult{
		Name:    rootPath,
		Subdirs: make(DirResults, 0),
		Files:   make([]FileResult, 0),
		Summary: make(map[string]int),
	}
	if filepath.Base(rootPath) == ".git" {
		logger.Printf("INFO Skipping %q.\n", rootPath)
		return result
	}
	// open(2) the directory to readdir(2) and stat(2) it.
	dir, err := os.Open(rootPath)
	if err != nil {
		logger.Println("ERROR", err)
		return result
	}
	defer dir.Close()
	fileinfoz, err := dir.Readdir(0)
	if err != nil {
		logger.Println("ERROR", err)
		return result
	}

	// Spawn one goroutine per subdirectory, and another one per file.
	dirResultsChan := make(chan DirResult)
	fileResultsChan := make(chan *FileResult)
	count := 0
	for _, fileinfo := range fileinfoz {
		filename := filepath.Join(rootPath, fileinfo.Name())
		if fileinfo.IsDir() {
			count++
			go func(path string) {
				dirResultsChan <- locDir(path)
			}(filename)
		} else if fileinfo.Mode().IsRegular() {
			count++
			go func(filename string) {
				fileResultsChan <- locFile(filename)
			}(filename)
		} else {
			logger.Printf("INFO Skipping non-regular and non-directory file %q.\n", filename)
		}
	}

	// Gather goroutines' results.
	for ; count > 0; count-- {
		select {
		case dr := <-dirResultsChan:
			result.Subdirs = append(result.Subdirs, dr)
			for lang, loc := range dr.Summary {
				if _, exists := result.Summary[lang]; exists {
					result.Summary[lang] += loc
				} else {
					result.Summary[lang] = loc
				}
			}
		case fr := <-fileResultsChan:
			if fr != nil {
				result.Files = append(result.Files, *fr)
				for lang, loc := range fr.Loc {
					if _, exists := result.Summary[lang]; exists {
						result.Summary[lang] += loc
					} else {
						result.Summary[lang] = loc
					}
				}
			}
		}
	}
	close(dirResultsChan)
	close(fileResultsChan)

	return result
}

// The core function for detecting a file's type, creating a LocCounter to
// count the lines of code in it, and finally return the results in a
// FileResult struct.
func locFile(filename string) *FileResult {
	var result *FileResult

	file, err := os.Open(filename)
	if err != nil {
		logger.Println("ERROR", err)
		return result
	}
	defer file.Close()

	baseName := filepath.Base(filename)
	ext := filepath.Ext(filename)
	if ext == "" {
		if strings.HasPrefix(baseName, "Makefile") {
			ext = "Makefile"
		} else if strings.HasPrefix(baseName, "Dockerfile") {
			ext = "Dockerfile"
		}
	} else {
		// Ignore the leading dot.
		ext = ext[1:]
	}
	locCounter, err := NewLocCounter(file, ext)
	if err != nil {
		logger.Println("ERROR", err)
		return result
	}

	loc, err := locCounter.Count()
	if err != nil {
		logger.Println("ERROR", err)
	}
	result = &FileResult{
		Name: baseName,
		Loc: map[string]int{
			languages[ext].name: loc,
		},
	}
	return result
}
