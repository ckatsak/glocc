/*
Copyright (C) 2017, Christos Katsakioris
All rights reserved.

This software may be modified and distributed under the terms
of the BSD 3-Clause License. See the LICENSE file for details.
*/

// A package implementing a counter of lines of code in files and directories.
//
// It also includes a command line tool, glocc, which can be used to easily
//
// glocc is massively parallel; counting each file and each subdirectory is
// assigned to a separate goroutine, and their separate results are later
// merged on a per-subdirectory basis.
//
// TODO: Add documentation for glocc usage as a package.
//
// ?TODO: Add documentation for the glocc command line tool?
package glocc
