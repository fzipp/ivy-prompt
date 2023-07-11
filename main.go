// Copyright 2022 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Ivy-prompt is a line editor interface wrapper for Rob Pike's [Ivy],
// an APL-like big number calculator. It provides a Readline-style input
// prompt with an input history and tab-completion.
//
// [Ivy]: https://robpike.io/ivy
package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/peterh/liner"
)

var (
	opMultiline   = regexp.MustCompile(`^\s*op\s*.*=\s*$`)
	specialDemo   = regexp.MustCompile(`^\s*\)\s*demo\s*$`)
	specialPrompt = regexp.MustCompile(`^\s*\)\s*prompt.*$`)
)

func main() {
	origMode, err := liner.TerminalMode()
	check(err)

	ivy, err := NewIvy()
	check(err)
	err = ivy.Start()
	check(err)
	defer printErr(ivy.Quit())

	fmt.Println(`Ivy big number calculator. https://robpike.io/ivy
Type ")help" for help, Ctrl-D to quit.`)

	l := liner.NewLiner()
	defer printErr(l.Close())
	l.SetWordCompleter(makeCompleter(ivy))
	l.SetTabCompletionStyle(liner.TabPrints)

	linerMode, err := liner.TerminalMode()
	if err != nil {
		printErr(fmt.Errorf("could not determine terminal mode: %w", err))
		return
	}

	_ = loadHistory(l)
	for err == nil {
		err = readEvalPrint(ivy, l, origMode, linerMode)
	}
	if err != io.EOF {
		printErr(err)
		// continue in order to save history
	}
	err = saveHistory(l)
	if err != nil {
		printErr(fmt.Errorf("could not save history: %w", err))
	}
}

// An io.EOF error is returned if the user signals end-of-file
// by pressing Ctrl-D.
func readEvalPrint(ivy *Ivy, l *liner.State, origMode, linerMode liner.ModeApplier) error {
	input, err := l.Prompt("ivy> ")
	if err != nil {
		return err
	}
	if specialDemo.MatchString(input) {
		err = origMode.ApplyMode()
		if err != nil {
			return fmt.Errorf("could not apply original terminal mode: %w", err)
		}
		runDemo()
		err = linerMode.ApplyMode()
		if err != nil {
			return fmt.Errorf("could not apply terminal mode for line editor: %w", err)
		}
	} else if specialPrompt.MatchString(input) {
		fmt.Println("prompt special command not yet supported in ivy-prompt")
	} else {
		if opMultiline.MatchString(input) {
			input, err = readMultiline(l, input)
			if err != nil {
				return err
			}
		}
		var resp string
		resp, err = ivy.Exec(input)
		if err != nil {
			return fmt.Errorf("could not execute input: %w", err)
		}
		fmt.Print(resp)
	}
	l.AppendHistory(input)
	return nil
}

func runDemo() {
	cmd := exec.Command("ivy", "-e", ")demo")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		printErr(fmt.Errorf("could not run demo: %w", err))
	}
}

// An io.EOF error is returned if the user signals end-of-file
// by pressing Ctrl-D.
func readMultiline(l *liner.State, firstLine string) (string, error) {
	var b strings.Builder
	b.WriteString(firstLine)
	for {
		line, err := l.Prompt("        ")
		if err != nil {
			return "", err
		}
		b.WriteByte('\n')
		b.WriteString(line)
		if strings.TrimSpace(line) == "" {
			break
		}
	}
	return b.String(), nil
}

func check(err error) {
	if err != nil {
		printErr(err)
		os.Exit(1)
	}
}

func printErr(err error) {
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
