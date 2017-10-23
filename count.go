/*
Copyright (C) 2017, Christos Katsakioris
All rights reserved.

This software may be modified and distributed under the terms
of the BSD 3-Clause License. See the LICENSE file for details.
*/

package glocc

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
	Name string         `json:"name" yaml:"Name"`
	Loc  map[string]int `json:"loc" yaml:"loc"`
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
		//if os.IsNotExist(err) {  // XXX What did I want this?
		logger.Println("ERROR", err)
		return result
	}
	if fileinfo.IsDir() {
		result = locDir(rootPath)
	} else if fileinfo.Mode().IsRegular() {
		fileResult := locFile(rootPath)
		result.Name = fileResult.Name
		result.Subdirs = nil
		result.Files = []FileResult{fileResult}
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
	fileResultsChan := make(chan FileResult)
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
			result.Files = append(result.Files, fr)
			for lang, loc := range fr.Loc {
				if _, exists := result.Summary[lang]; exists {
					result.Summary[lang] += loc
				} else {
					result.Summary[lang] = loc
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
func locFile(filename string) FileResult {
	result := FileResult{Loc: make(map[string]int)}

	file, err := os.Open(filename)
	if err != nil {
		logger.Println("ERROR", err)
		return result
	}
	defer file.Close()

	ext := filepath.Ext(filename)
	if ext == "" {
		baseName := filepath.Base(filename)
		if baseName == "Makefile" {
			ext = "Makefile"
		} else if baseName == "Dockerfile" {
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

	if loc, err := locCounter.Count(); err != nil {
		logger.Println("ERROR", err)
	} else {
		result.Loc[languages[ext].name] = loc
		result.Name = filepath.Base(filename)
	}
	return result
}
