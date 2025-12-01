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
	specialHelp  = regexp.MustCompile(`^\s*\)\s*help\s*$`)
	specialDebug = regexp.MustCompile(`^\s*\)\s*debug\s*$`)
	specialClear = regexp.MustCompile(`^\s*\)\s*clear\s*$`)
	specialGet   = regexp.MustCompile(`^\s*\)\s*get\s*["']?$`)
	specialSave  = regexp.MustCompile(`^\s*\)\s*save\s*["']?$`)
	opSys        = regexp.MustCompile(`(|.*\s+)sys(|\s*["'])$`)
)

type defsProvider interface {
	Ops() ([]string, error)
	Vars() ([]string, error)
}

func makeCompleter(def defsProvider) liner.WordCompleter {
	return func(line string, pos int) (head string, completions []string, tail string) {
		idx := 0
		for range pos {
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
		} else if specialDebug.MatchString(beforeWord) {
			words = debugSettings
		} else if specialGet.MatchString(beforeWord) || specialSave.MatchString(beforeWord) {
			words = fileNames()
			suffix = ""
		} else if specialClear.MatchString(beforeWord) {
			words = clearTopics
		} else if opSys.MatchString(beforeWord) {
			last := beforeWord[len(beforeWord)-1]
			if last == '"' || last == '\'' {
				words = sysTopics
			} else {
				words = []string{"\""}
			}
			suffix = ""
		} else {
			ops, err := def.Ops()
			if err == nil {
				words = append(words, ops...)
			}
			vars, err := def.Vars()
			if err == nil {
				words = append(words, vars...)
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
		// don't report error, because file name completion
		// is optional
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
