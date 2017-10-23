# glocc

[![GoDoc](https://godoc.org/github.com/ckatsak/glocc?status.png)](https://godoc.org/github.com/ckatsak/glocc)
[![Go Report Card](https://goreportcard.com/badge/github.com/ckatsak/glocc)](https://goreportcard.com/badge/github.com/ckatsak/glocc)

`glocc` is a package implementing a relatively fast, parallel counter of lines
of code in files and directories.

It also includes a command line tool, `glocc`, which is handy for performing
such counting and pretty (brief or extensive) printing of the results.

`glocc` is an **aggressively parallel** solution to an embarrassingly parallel
problem. The count of every file and every subdirectory is assigned to a
separate goroutine. All spawned goroutines are properly synchronized and their
independent results are merged later, on a higher level (level = on a
per-subdirectory basis).

It was originally written for use with personal projects and small codebases,
and also to get in touch with the Go programming language.
Performance-wise, it can be further improved (and hopefully will be, when I
have more time).

## Contents

- [Command line tool](#command-line-tool)
- [Installation](#installation)
- [Platforms](#platforms)
- [Supported Languages](#supported-languages)
- [`glocc` as package](#glocc-as-package)
- [Known Issue](#known-issues)

## Command line tool <a name="command-line-tool"></a>

Simply run it with any number of files or directories as command line
arguments:
```text
$ glocc ~/foo src/bar
```

By default, only a summary of all counted lines is printed to the standard
output. To print the results extensively in a tree-like format, it can be
executed with the `-a` flag:
```text
$ glocc -a baz.go ~/src/foo
```

The results can be printed in **YAML** (default) or **JSON** format, using the
`-o` flag:
```text
$ glocc -o json ~/bar
```

Running it with the `-h` flag shows all options available.

## Installation <a name="installation"></a>

For both the package and the command line tool to be installed, assuming [Go is
properly installed](https://golang.org/doc/install), it should be as easy as:
```text
$ go get -u github.com/ckatsak/glocc/...
```

## Platforms <a name="platforms"></a>

Until now, it has been tested only under `go version go1.9.1 linux/amd64`.

## Supported Languages <a name="supported-languages"></a>

- Ada
- Assembly
- AWK
- C
- C++
- C#
- D (not the ddoc comments)
- Delphi
- Dockerfile
- Eiffel
- Elixir
- Erlang
- Go
- Haskell
- HTML
- Java
- Javascript
- JSON
- Kotlin
- Lisp
- Makefile
- Matlab
- OCaml
- Perl (not `__END__` comments)
- PHP
- PowerShell
- Python
- R
- Ruby (not `__END__` comments)
- Rust
- Scala
- Scheme
- shell scripts
- SQL
- Standard ML
- TeX
- Tcl
- YAML

## Using the `glocc` package <a name="glocc-as-package"></a>

For use as a package, `glocc` exports `func CountLoc(root string) DirResult`,
which, given a root directory, returns a struct of type `DirResult`, a
custom (recursive) type that contains the results of counting all lines of
code under this root directory.

It also exports `EnableLogging()` and `DisableLogging()` functions, to enable
and disable verbose logging to standard error, respectively, using a
package-level logger.
Note that verbose logging includes details about every line of every file
visited, which might be quite ...verbose, and not that useful.

## Known Issue <a name="known-issues"></a>

For now, really huge source trees, like the Linux kernel source tree, might
rarely cause `glocc` to crash, due the big number of blocked OS threads trying
to handle the huge number of goroutines spawned. To be more precise, the exact
problem is reported as:

```text
$ glocc ./linux
runtime: program exceeds 10000-thread limit
fatal error: thread exhaustion
```

It *cannot* occur in small and medium-sized codebases, and it's also unlikely
to occur in bigger ones too. Just be warned.
I plan to hack around this problem once I have the time; maybe using some kind
of pool or something, or by spawning the goroutines in some clever way.
As long as this note is here though, the bug is probably still around.
Theoretically, a quick and dirty solution would be to increase the number of
operating system threads that a Go program can use, using the `SetMaxThreads()`
function in package `runtime/debug`; the default value is set to 10000 threads.
However, mind that (quoted from the
[official documentation](https://golang.org/pkg/runtime/debug/#SetMaxThreads)):

> SetMaxThreads is useful mainly for limiting the damage done by programs
> that create an unbounded number of threads. The idea is to take down
> the program before it takes down the operating system.
