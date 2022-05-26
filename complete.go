// Copyright 2022 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/peterh/liner"
)

var specialHelp = regexp.MustCompile(`^\s*\)\s*help\s*$`)

func makeCompleter(ivy *Ivy) liner.WordCompleter {
	return func(line string, pos int) (head string, completions []string, tail string) {
		idx := 0
		for i := 0; i < pos; i++ {
			_, width := utf8.DecodeRuneInString(line[idx:])
			idx += width
		}
		wordStart := strings.LastIndexFunc(line[:idx], isNotAlphaNum) + 1
		partialWord := strings.ToLower(line[wordStart:idx])
		beforeWord := strings.TrimSpace(line[:wordStart])
		words := keywords
		if beforeWord == ")" {
			words = specials
		} else if specialHelp.MatchString(beforeWord) {
			words = append(words, helpTopics...)
		} else {
			ops, err := ivy.Ops()
			if err == nil {
				words = append(words, ops...)
			}
		}
		for _, word := range words {
			if strings.HasPrefix(strings.ToLower(word), partialWord) {
				completions = append(completions, word+" ")
			}
		}
		sort.Strings(completions)
		return line[:wordStart], completions, line[idx:]
	}
}

func isNotAlphaNum(r rune) bool {
	return !unicode.IsLetter(r) && !unicode.IsDigit(r)
}
