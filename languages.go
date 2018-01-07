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

// A struct to store all the basic information needed to support counting the
// lines of code for a programming language, hardcoded.
type language struct {
	name       string
	extensions []string

	inlineCommentTokens            []string
	multiLineCommentStartingTokens []string
	multiLineCommentEndingTokens   []string
}

// A slice of language structs containing all the programming languages
// currently supported by glocc.
var allLanguages = []language{
	{
		name:                           "Ada",
		extensions:                     []string{"adb", "ads"},
		inlineCommentTokens:            []string{`--`},
		multiLineCommentStartingTokens: []string{},
		multiLineCommentEndingTokens:   []string{},
	},
	{
		name:                           "Assembly",
		extensions:                     []string{"asm", "s", "S"},
		inlineCommentTokens:            []string{`;`}, // works for NASM, but not for every assembly out there
		multiLineCommentStartingTokens: []string{},
		multiLineCommentEndingTokens:   []string{},
	},
	{
		name:                           "AWK",
		extensions:                     []string{"awk"},
		inlineCommentTokens:            []string{`#`},
		multiLineCommentStartingTokens: []string{},
		multiLineCommentEndingTokens:   []string{},
	},
	{
		name:                           "C",
		extensions:                     []string{"c", "h"},
		inlineCommentTokens:            []string{`//`},
		multiLineCommentStartingTokens: []string{`/*`},
		multiLineCommentEndingTokens:   []string{`*/`},
	},
	{
		name:                           "C++",
		extensions:                     []string{"cc", "hh", "C", "H", "cpp", "hpp", "cxx", "hxx", "c++", "h++"},
		inlineCommentTokens:            []string{`//`},
		multiLineCommentStartingTokens: []string{`/*`},
		multiLineCommentEndingTokens:   []string{`*/`},
	},
	{
		name:                           "C#",
		extensions:                     []string{"cs"},
		inlineCommentTokens:            []string{`//`, `///`},
		multiLineCommentStartingTokens: []string{`/*`, `/**`},
		multiLineCommentEndingTokens:   []string{`*/`},
	},
	{
		name:                           "D",
		extensions:                     []string{"d"},
		inlineCommentTokens:            []string{`//`, `///`},
		multiLineCommentStartingTokens: []string{`/*`, `/+`}, // nesting is supported, missing ddoc comment tokens
		multiLineCommentEndingTokens:   []string{`*/`, `+/`}, // nesting is supported
	},
	{
		name:                           "Delphi",
		extensions:                     []string{"p", "pp", "pas"},
		inlineCommentTokens:            []string{`//`},
		multiLineCommentStartingTokens: []string{`(*`, `{`},
		multiLineCommentEndingTokens:   []string{`*)`, `}`},
	},
	{
		name:                           "Dockerfile",
		extensions:                     []string{"Dockerfile"},
		inlineCommentTokens:            []string{`#`},
		multiLineCommentStartingTokens: []string{},
		multiLineCommentEndingTokens:   []string{},
	},
	{
		name:                           "Eiffel",
		extensions:                     []string{"e"},
		inlineCommentTokens:            []string{`--`},
		multiLineCommentStartingTokens: []string{},
		multiLineCommentEndingTokens:   []string{},
	},
	{
		name:                           "Elixir",
		extensions:                     []string{"ex", "exs"},
		inlineCommentTokens:            []string{`%`},
		multiLineCommentStartingTokens: []string{},
		multiLineCommentEndingTokens:   []string{},
	},
	{
		name:                           "Erlang",
		extensions:                     []string{"erl", "hrl"},
		inlineCommentTokens:            []string{`%`},
		multiLineCommentStartingTokens: []string{},
		multiLineCommentEndingTokens:   []string{},
	},
	{
		name:                           "Go",
		extensions:                     []string{"go"},
		inlineCommentTokens:            []string{`//`},
		multiLineCommentStartingTokens: []string{`/*`},
		multiLineCommentEndingTokens:   []string{`*/`},
	},
	{
		name:                           "Haskell",
		extensions:                     []string{"hs", "lhs"},
		inlineCommentTokens:            []string{`--`},
		multiLineCommentStartingTokens: []string{`{-`}, // nesting is not supported
		multiLineCommentEndingTokens:   []string{`-}`}, // nesting is not supported
	},
	{
		name:                           "HTML",
		extensions:                     []string{"html", "htm"},
		inlineCommentTokens:            []string{},
		multiLineCommentStartingTokens: []string{`<!--`},
		multiLineCommentEndingTokens:   []string{`-->`},
	},
	{
		name:                           "Java",
		extensions:                     []string{"java"},
		inlineCommentTokens:            []string{`//`},
		multiLineCommentStartingTokens: []string{`/*`, `/**`},
		multiLineCommentEndingTokens:   []string{`*/`},
	},
	{
		name:                           "Javascript",
		extensions:                     []string{"js"},
		inlineCommentTokens:            []string{`//`},
		multiLineCommentStartingTokens: []string{`/*`},
		multiLineCommentEndingTokens:   []string{`*/`},
	},
	{
		name:                           "JSON",
		extensions:                     []string{"json"},
		inlineCommentTokens:            []string{},
		multiLineCommentStartingTokens: []string{},
		multiLineCommentEndingTokens:   []string{},
	},
	{
		name:                           "Kotlin",
		extensions:                     []string{"kt", "kts"},
		inlineCommentTokens:            []string{`//`},
		multiLineCommentStartingTokens: []string{`/*`},
		multiLineCommentEndingTokens:   []string{`*/`},
	},
	{
		name:                           "Lisp",
		extensions:                     []string{"lisp", "lsp", "l", "cl", "fasl"},
		inlineCommentTokens:            []string{`;`},
		multiLineCommentStartingTokens: []string{`#|`},
		multiLineCommentEndingTokens:   []string{`|#`},
	},
	{
		name:                           "Makefile",
		extensions:                     []string{"Makefile"},
		inlineCommentTokens:            []string{`#`},
		multiLineCommentStartingTokens: []string{},
		multiLineCommentEndingTokens:   []string{},
	},
	{
		name:                           "Markdown",
		extensions:                     []string{"md"},
		inlineCommentTokens:            []string{},
		multiLineCommentStartingTokens: []string{},
		multiLineCommentEndingTokens:   []string{},
	},
	{
		name:                           "Matlab",
		extensions:                     []string{"m"},
		inlineCommentTokens:            []string{`%`},
		multiLineCommentStartingTokens: []string{`%{`},
		multiLineCommentEndingTokens:   []string{`%}`},
	},
	{
		name:                           "OCaml",
		extensions:                     []string{"ml", "mli"},
		inlineCommentTokens:            []string{},
		multiLineCommentStartingTokens: []string{`(*`}, // nesting is not supported
		multiLineCommentEndingTokens:   []string{`*)`}, // nesting is not supported
	},
	{
		name:                           "Perl",
		extensions:                     []string{"pl", "pm", "t", "pod"},
		inlineCommentTokens:            []string{`#`},
		multiLineCommentStartingTokens: []string{`=begin`}, // __END__ is not supported
		multiLineCommentEndingTokens:   []string{`=cut`},
	},
	{
		name:                           "PHP",
		extensions:                     []string{"php"},
		inlineCommentTokens:            []string{`#`, `//`},
		multiLineCommentStartingTokens: []string{`/*`, `/**`},
		multiLineCommentEndingTokens:   []string{`*/`},
	},
	{
		name:                           "PowerShell",
		extensions:                     []string{"ps1"},
		inlineCommentTokens:            []string{`#`},
		multiLineCommentStartingTokens: []string{`<#`},
		multiLineCommentEndingTokens:   []string{`#>`},
	},
	{
		name:                           "Protocol Buffers",
		extensions:                     []string{"proto"},
		inlineCommentTokens:            []string{`//`},
		multiLineCommentStartingTokens: []string{`/*`},
		multiLineCommentEndingTokens:   []string{`*/`},
	},
	{
		name:                           "Python",
		extensions:                     []string{"py"},
		inlineCommentTokens:            []string{`#`},
		multiLineCommentStartingTokens: []string{`"""`, `'''`}, // nesting is supported
		multiLineCommentEndingTokens:   []string{`"""`, `'''`}, // nesting is supported
	},
	{
		name:                           "R",
		extensions:                     []string{"r", "R", "RData", "rds", "rda"},
		inlineCommentTokens:            []string{`#`},
		multiLineCommentStartingTokens: []string{},
		multiLineCommentEndingTokens:   []string{},
	},
	{
		name:                           "Ruby",
		extensions:                     []string{"rb"},
		inlineCommentTokens:            []string{`#`},
		multiLineCommentStartingTokens: []string{`=begin`}, // __END__ is not supported
		multiLineCommentEndingTokens:   []string{`=end`},
	},
	{
		name:                           "Rust",
		extensions:                     []string{"rs", "rlib"},
		inlineCommentTokens:            []string{`//`, `///`, `//!`},
		multiLineCommentStartingTokens: []string{`/*`, `/**`, `/*!`},
		multiLineCommentEndingTokens:   []string{`*/`},
	},
	{
		name:                           "Scala",
		extensions:                     []string{"scala", "sc"},
		inlineCommentTokens:            []string{`//`},
		multiLineCommentStartingTokens: []string{`/*`},
		multiLineCommentEndingTokens:   []string{`*/`},
	},
	{
		name:                           "Scheme",
		extensions:                     []string{"scm", "ss"},
		inlineCommentTokens:            []string{`;`},
		multiLineCommentStartingTokens: []string{`#|`},
		multiLineCommentEndingTokens:   []string{`|#`},
	},
	{
		name:                           "Shell",
		extensions:                     []string{"sh", "bash", "zsh", "ksh", "csh"},
		inlineCommentTokens:            []string{`#`},
		multiLineCommentStartingTokens: []string{},
		multiLineCommentEndingTokens:   []string{},
	},
	{
		name:                           "SQL",
		extensions:                     []string{"sql"},
		inlineCommentTokens:            []string{`--`},
		multiLineCommentStartingTokens: []string{},
		multiLineCommentEndingTokens:   []string{},
	},
	{
		name:                           "Standard ML",
		extensions:                     []string{"sml"},
		inlineCommentTokens:            []string{},
		multiLineCommentStartingTokens: []string{`(*`},
		multiLineCommentEndingTokens:   []string{`*)`},
	},
	{
		name:                           "TeX",
		extensions:                     []string{"tex"},
		inlineCommentTokens:            []string{`%`},
		multiLineCommentStartingTokens: []string{},
		multiLineCommentEndingTokens:   []string{},
	},
	{
		name:                           "plain text",
		extensions:                     []string{"txt"},
		inlineCommentTokens:            []string{},
		multiLineCommentStartingTokens: []string{},
		multiLineCommentEndingTokens:   []string{},
	},
	{
		name:                           "Tcl",
		extensions:                     []string{"tcl", "tbc"},
		inlineCommentTokens:            []string{`#`},
		multiLineCommentStartingTokens: []string{},
		multiLineCommentEndingTokens:   []string{},
	},
	{
		name:                           "YAML",
		extensions:                     []string{"yaml", "yml"},
		inlineCommentTokens:            []string{`#`},
		multiLineCommentStartingTokens: []string{},
		multiLineCommentEndingTokens:   []string{},
	},
}

// Map file extensions to language structs, for fast looking up.
var languages = map[string]language{}

func init() {
	// Populate global var languages.
	for _, lang := range allLanguages {
		for _, ext := range lang.extensions {
			languages[ext] = lang
		}
	}
}
