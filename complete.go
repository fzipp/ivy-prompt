// Copyright 2022 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/peterh/liner"
)

var (
	specialHelp = regexp.MustCompile(`^\s*\)\s*help\s*$`)
	specialGet  = regexp.MustCompile(`^\s*\)\s*get\s*["']$`)
	specialSave = regexp.MustCompile(`^\s*\)\s*save\s*["']$`)
)

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
		suffix := " "
		words := keywords
		if beforeWord == ")" {
			words = specials
		} else if specialHelp.MatchString(beforeWord) {
			words = append(words, helpTopics...)
		} else if specialGet.MatchString(beforeWord) || specialSave.MatchString(beforeWord) {
			words = fileNames()
			suffix = ""
		} else {
			ops, err := ivy.Ops()
			if err == nil {
				words = append(words, ops...)
			}
		}
		for _, word := range words {
			if strings.HasPrefix(strings.ToLower(word), partialWord) {
				completions = append(completions, word+suffix)
			}
		}
		sort.Strings(completions)
		return line[:wordStart], completions, line[idx:]
	}
}

func fileNames() []string {
	entries, err := os.ReadDir(".")
	if err != nil {
		return nil
	}
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files
}

func isNotAlphaNum(r rune) bool {
	return !unicode.IsLetter(r) && !unicode.IsDigit(r)
}
