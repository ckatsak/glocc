/*
Copyright (C) 2017, Christos Katsakioris
All rights reserved.

This software may be modified and distributed under the terms
of the BSD 3-Clause License. See the LICENSE file for details.
*/

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
		name:                           "Assembly",
		extensions:                     []string{"asm", "s", "S"},
		inlineCommentTokens:            []string{`;`}, // works for NASM, but not for every assembly out there
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
		name:                           "D",
		extensions:                     []string{"d"},
		inlineCommentTokens:            []string{`//`, `///`},
		multiLineCommentStartingTokens: []string{`/*`, `/+`}, // nesting is supported, missing some tokens here
		multiLineCommentEndingTokens:   []string{`*/`, `+/`}, // nesting is supported
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
		multiLineCommentStartingTokens: []string{`{-`}, // nesting is supported
		multiLineCommentEndingTokens:   []string{`-}`}, // nesting is supported
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
		name:                           "Kotlin",
		extensions:                     []string{"kt", "kts"},
		inlineCommentTokens:            []string{`//`},
		multiLineCommentStartingTokens: []string{`/*`},
		multiLineCommentEndingTokens:   []string{`*/`},
	},
	{
		name:                           "Makefile",
		extensions:                     []string{"Makefile"},
		inlineCommentTokens:            []string{`#`},
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
		multiLineCommentStartingTokens: []string{`(*`}, // nesting is supported
		multiLineCommentEndingTokens:   []string{`*)`}, // nesting is supported
	},
	{
		name:                           "PHP",
		extensions:                     []string{"php"},
		inlineCommentTokens:            []string{`#`, `//`},
		multiLineCommentStartingTokens: []string{`/*`, `/**`},
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
		name:                           "Ruby",
		extensions:                     []string{"rb"},
		inlineCommentTokens:            []string{`#`},
		multiLineCommentStartingTokens: []string{`=begin`},
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
		name:                           "Shell",
		extensions:                     []string{"sh", "bash", "zsh", "ksh", "csh"},
		inlineCommentTokens:            []string{`#`},
		multiLineCommentStartingTokens: []string{},
		multiLineCommentEndingTokens:   []string{},
	},
	{
		name:                           "SML",
		extensions:                     []string{"sml"},
		inlineCommentTokens:            []string{},
		multiLineCommentStartingTokens: []string{`(*`},
		multiLineCommentEndingTokens:   []string{`*)`},
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
