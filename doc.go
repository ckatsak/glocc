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
