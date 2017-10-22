/*
Copyright (C) 2017, Christos Katsakioris
All rights reserved.

This software may be modified and distributed under the terms
of the BSD 3-Clause License. See the LICENSE file for details.
*/
package glocc

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

// The results of counting lines of code in every file and subdirectory of a
// directory are gathered in a DirResult.
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

// A slice of DirResult.
type DirResults []DirResult

// The results of counting lines of code in a file are stored in a FileResult,
// which is eventually placed in the DirResult associated with the directory
// that the file lives under.
type FileResult struct {
	Name string         `json:"name" yaml:"Name"`
	Loc  map[string]int `json:"loc" yaml:"loc"`
}

// Package-level logger.
var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "glocc: ", log.Ltime|log.Lmicroseconds|log.Lshortfile)
}

// This function is the exported interface of glocc package, meant to be called
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
		logger.Println(err)
		return result
	}
	fileinfo, err := os.Stat(rootPath)
	if os.IsNotExist(err) {
		logger.Println(err)
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
	logger.Printf("Time elapsed for %q: %s\n", root, time.Since(start))
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
		logger.Printf("Skipping %q.\n", rootPath)
		return result
	}
	// open(2) the directory to readdir(2) and stat(2) it.
	dir, err := os.Open(rootPath)
	if err != nil {
		logger.Println(err)
		return result
	}
	defer dir.Close()
	fileinfoz, err := dir.Readdir(0)
	if err != nil {
		logger.Println(err)
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
			logger.Printf("Skipping non-regular and non-directory file %q.\n", filename)
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
		logger.Println(err)
		return result
	}
	defer file.Close()

	ext := filepath.Ext(filename)
	if ext == "" {
		if filepath.Base(filename) == "Makefile" {
			ext = "Makefile"
		}
	} else {
		// Ignore the leading dot.
		ext = ext[1:]
	}
	locCounter, err := NewLocCounter(file, ext)
	if err != nil {
		logger.Println(err)
		return result
	}

	if loc, err := locCounter.Count(); err != nil {
		logger.Println(err)
	} else {
		result.Loc[languages[ext].name] = loc
		result.Name = filepath.Base(filename)
	}
	return result
}
