/*
Copyright (C) 2017, Christos Katsakioris
All rights reserved.

This software may be modified and distributed under the terms
of the BSD 3-Clause License. See the LICENSE file for details.
*/

// Package glocc implements a very fast, parallel counter of lines of code in
// files and directories.
//
// It also includes a command line tool, glocc, which is extremely fast and
// handy for performing such counting and pretty printing, either briefly or
// extensively, the results.
//
// glocc is aggressively (...embarrassingly) parallel. The counting for every
// file and every subdirectory is assigned to a separate goroutine. All
// goroutines spawned are properly synchronized and their independent results
// are merged later, on a higher level (level = on a per-subdirectory basis).
//
// glocc command line tool
//
// Simply call it with any number of files or directories as command line
// arguments:
//
//	$ glocc ~/project1 ~/project2
//
// By default, only a summary of all counted lines is printed to the standard
// output. To print the results extensively in a tree-like structure, it can be
// executed with flag:
//
//	$ glocc -a  ~/project1 ~/project2
//
// glocc as a package
//
// For use as a package, glocc exports `func CountLoc(root string) DirResult`,
// which, given a root directory, returns a struct of type `DirResult`, a
// custom (recursive) type that contains the results of counting all lines of
// code under this root directory.
//
// It also exports EnableLogging() and DisableLogging() functions, to enable
// and disable verbose logging to standard error, respectively, using a
// package-level logger.
// Note that verbose logging includes details about every line of every file
// visited, which might be quite ...verbose, and not that useful.
//
// Known bug
//
// For now, really huge source trees, like the Linux kernel source tree, might
// crash glocc, due the big number of blocked OS threads trying to handle the
// huge number of goroutines spawned. To be more precise, the problem is:
//
//	$ glocc ./linux
// 	runtime: program exceeds 10000-thread limit
// 	fatal error: thread exhaustion
//
// It is planned to hack around this problem, maybe using some kind of pool or
// something; as long as this note is here though, the bug is probably still
// around.
// Theoretically, a quick and dirty solution would be to increase the number of
// operating system threads that a Go program can use, using SetMaxThreads() of
// runtime/debug; the default value is set to 10000 threads. However, mind that
// (quoted from https://golang.org/pkg/runtime/debug/#SetMaxThreads):
//
// 	SetMaxThreads is useful mainly for limiting the damage done by programs
//	that create an unbounded number of threads. The idea is to take down
//	the program before it takes down the operating system.
package glocc
