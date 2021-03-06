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
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

// These states don't need to exist per LocCounter, as they don't carry any
// LocCounter-specific data.
var (
	globalStateInitial = &stateInitial{}
	globalStateCode    = &stateCode{}
)

// LocCounter is the core entity of the package, which initiates and later
// holds the state of the counting for a single file.
// It is associated to the counting of a single file, and created in the
// goroutine that is assigned to count the file.
type LocCounter struct {
	language language
	loc      int

	file            *os.File
	currLine        string
	currLineCounted bool
	fileLinesCnt    int

	state                 loccState
	stateMultiLineComment *stateMultiLineComment
}

// NewLocCounter returns a new LocCounter, properly initialized to count the
// lines of code in a specific file of a specific language.
// Returns an error if a supported language cannot be detected.
func NewLocCounter(file *os.File, ext string) (lc *LocCounter, err error) {
	if lang, valid := languages[ext]; !valid {
		err = fmt.Errorf("Cannot deduce a supported language from extension %q.", ext)
	} else {
		lc = &LocCounter{
			language: lang,
			file:     file,
			state:    globalStateInitial,
			stateMultiLineComment: &stateMultiLineComment{},
		}
	}
	return
}

// Count is the only exported method of LocCounter. It basically reads (line by
// line) the content of the file associated with the LocCounter, and performs
// the counting. It is implemented using the State design pattern.
func (lc *LocCounter) Count() (int, error) {
	logger.Printf("DEBUG LocCounter.Count() for file %q: Starting...\n", lc.file.Name())
	fsc := bufio.NewScanner(lc.file)
	for fsc.Scan() {
		lc.fileLinesCnt++
		lc.currLine = fsc.Text()
		lc.currLine = strings.TrimLeft(lc.currLine, " \t") // trim leading whitespace
		lc.currLineCounted = false
		for !lc.state.process(lc) {
		}
		if lc.currLineCounted {
			logger.Printf("DEBUG %q:%d --> Counted\n", lc.file.Name(), lc.fileLinesCnt)
			lc.loc++
		} else {
			logger.Printf("DEBUG %q:%d --> Discarded\n", lc.file.Name(), lc.fileLinesCnt)
		}
	}
	if err := fsc.Err(); err != nil {
		logger.Println("ERROR", err)
		return lc.loc, err
	}

	logger.Printf("DEBUG LocCounter.Count() for file %q: Finished.\n", lc.file.Name())
	return lc.loc, nil
}

// Change the state of the LocCounter.
func (lc *LocCounter) setState(state loccState) {
	lc.state = state
}

// Returns true if current line is empty; false otherwise.
func (lc *LocCounter) lineIsEmpty() bool {
	if len(lc.currLine) == 0 {
		return true
	}
	return false
}

// Returns the index of the first inline comment token that was found in
// current line, or the length of current line if none was found.
func (lc *LocCounter) inlineCommentIndex() int {
	firstInlineCommTokenIdx := len(lc.currLine)
	for _, t := range lc.language.inlineCommentTokens {
		ilcIdx := strings.Index(lc.currLine, t)
		if ilcIdx != -1 && ilcIdx < firstInlineCommTokenIdx {
			firstInlineCommTokenIdx = ilcIdx
		}
	}
	if firstInlineCommTokenIdx < len(lc.currLine) {
		logger.Printf("DEBUG Inline comment token found at %q:%d\n", lc.file.Name(), lc.fileLinesCnt)
	}
	return firstInlineCommTokenIdx
}

// The current state of a LocCounter. It may change from zero to multiple times
// while processing the same single line.
// Part of the State design pattern implementation.
type loccState interface {
	// The bool returned shows whether we're done processing currLine, so
	// as to break from the loop that LoccState.process() was called in.
	process(*LocCounter) bool
}

// The initial state in which every LocCounter starts in.
type stateInitial struct{}

// Line processing method for state stateInitial.
func (s *stateInitial) process(lc *LocCounter) bool {
	firstInlineCommTokenIdx := lc.inlineCommentIndex()
	if lc.lineIsEmpty() || firstInlineCommTokenIdx == 0 {
		return true
	}
	// On the first non-empty and non-inline-commented-out line, the state is changing.
	// Find the first occurrence of a multi-line comment starting token, if any.
	firstMultiLineCommTokenIdx, firstMultiLineCommToken := len(lc.currLine), ""
	for _, t := range lc.language.multiLineCommentStartingTokens {
		mlcIdx := strings.Index(lc.currLine, t)
		if mlcIdx != -1 && mlcIdx < firstMultiLineCommTokenIdx {
			firstMultiLineCommTokenIdx = mlcIdx
			firstMultiLineCommToken = t
		}
	}
	// If a multi-line comment starting token was found before the first inline comment token
	if firstMultiLineCommTokenIdx < firstInlineCommTokenIdx {
		logger.Printf("DEBUG Multi-line comment starting at %q:%d\n", lc.file.Name(), lc.fileLinesCnt)
		// If it wasn't in the beginning of the line
		if firstMultiLineCommTokenIdx > 0 {
			lc.currLineCounted = true
		}
		// Immediately continue processing the rest of the line in stateMultiLineComment,
		// as the state may change again within the same line.
		lc.currLine = strings.TrimLeft(lc.currLine[(firstMultiLineCommTokenIdx+len(firstMultiLineCommToken)):], " \t")
		lc.stateMultiLineComment.setToken(firstMultiLineCommToken)
		lc.setState(lc.stateMultiLineComment)
	} else {
		// If no multi-line comment starting token was found before the first inline comment token
		lc.setState(globalStateCode)
	}
	// State has to change from stateInitial in any case.
	return false
}

// The state of the LocCounter currently processing multi-line commented code.
type stateMultiLineComment struct {
	// Needed for Python (or any other language that I may not know of,
	// similar to Python in) that they need to nest e.g. occurrences of
	// `'''` in a `"""` multi-line comment, and of `"""` in a `'''`
	// multi-line comment.
	token string
}

// Line processing method for state stateMultiLineComment.
func (s *stateMultiLineComment) process(lc *LocCounter) bool {
	// Based on the observation that all supported languages actually use the
	// same token for closing block comments as for opening, only reversed.
	// Exceptions (handled) to this (for now): Ruby, and Java, PHP for docstrings.
	tokens := []string{} // the tokens which change the state
	reversedToken := reversed(lc.stateMultiLineComment.token)
	reversedTokenIsValid := false
	for _, t := range lc.language.multiLineCommentEndingTokens {
		if t == reversedToken {
			reversedTokenIsValid = true
			break
		}
	}
	if reversedTokenIsValid {
		tokens = append(tokens, reversedToken)
	} else {
		tokens = append(tokens, lc.language.multiLineCommentEndingTokens...)
	}

	// Find the first occurrence of a multi-line comment ending token, if any
	firstMultiLineCommTokenIdx, firstMultiLineCommToken := len(lc.currLine), ""
	for _, t := range tokens {
		mlcIdx := strings.Index(lc.currLine, t)
		if mlcIdx != -1 && mlcIdx < firstMultiLineCommTokenIdx {
			firstMultiLineCommTokenIdx = mlcIdx
			firstMultiLineCommToken = t
		}
	}
	// If a multi-line comment ending token was found
	if firstMultiLineCommTokenIdx < len(lc.currLine) {
		logger.Printf("DEBUG Multi-line comment ending at %q:%d\n", lc.file.Name(), lc.fileLinesCnt)
		s.token = ""
		lc.currLine = strings.TrimLeft(lc.currLine[(firstMultiLineCommTokenIdx+len(firstMultiLineCommToken)):], " \t")
		lc.setState(globalStateCode)
		return false
	}
	// If no multi-line comment ending token was found
	return true
}

// Change the saved token in stateMultiLineComment, and return the state struct
// itself.
func (s *stateMultiLineComment) setToken(token string) {
	s.token = token
}

// The state of the LocCounter currently processing code that needs to be
// counted in.
type stateCode struct{}

// Line processing method for state stateCode.
func (s *stateCode) process(lc *LocCounter) bool {
	firstInlineCommTokenIdx := lc.inlineCommentIndex()
	if lc.lineIsEmpty() || firstInlineCommTokenIdx == 0 {
		return true
	}
	// Find the first occurrence of a multi-line comment starting token, if any.
	firstMultiLineCommTokenIdx, firstMultiLineCommToken := len(lc.currLine), ""
	for _, t := range lc.language.multiLineCommentStartingTokens {
		mlcIdx := strings.Index(lc.currLine, t)
		if mlcIdx != -1 && mlcIdx < firstMultiLineCommTokenIdx {
			firstMultiLineCommTokenIdx = mlcIdx
			firstMultiLineCommToken = t
		}
	}
	// If a multi-line comment starting token was found before the first occurrence of an inline comment token
	if firstMultiLineCommTokenIdx < firstInlineCommTokenIdx {
		logger.Printf("DEBUG Multi-line comment start found at %q:%d\n", lc.file.Name(), lc.fileLinesCnt)
		// If it wasn't in the beginning of the line
		if firstMultiLineCommTokenIdx > 0 {
			lc.currLineCounted = true
		}
		// Immediately continue processing the rest of the line in stateMultiLineComment,
		// as the state may change again within the same line.
		lc.currLine = strings.TrimLeft(lc.currLine[(firstMultiLineCommTokenIdx+len(firstMultiLineCommToken)):], " \t")
		lc.stateMultiLineComment.setToken(firstMultiLineCommToken)
		lc.setState(lc.stateMultiLineComment)
		return false
	}
	lc.currLineCounted = true
	return true
}

// Returns the input string reversed.
func reversed(s string) string {
	size := len(s)
	buf := make([]byte, size)
	for i := 0; i < size; {
		r, n := utf8.DecodeRuneInString(s[i:])
		i += n
		utf8.EncodeRune(buf[(size-i):], r)
	}
	return string(buf)
}
