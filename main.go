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
	defer ivy.Quit()

	fmt.Println("Ivy big number calculator. https://robpike.io/ivy\nType \")help\" for help, Ctrl-D to quit.")

	l := liner.NewLiner()
	defer l.Close()
	l.SetWordCompleter(makeCompleter(ivy))
	l.SetTabCompletionStyle(liner.TabPrints)

	linerMode, err := liner.TerminalMode()
	if err != nil {
		return
	}

	_ = loadHistory(l)
	for err == nil {
		err = readEvalPrint(ivy, l, origMode, linerMode)
	}
	if err != io.EOF {
		fmt.Fprintln(os.Stderr, err)
	}
	err = saveHistory(l)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not save history: %s\n", err)
	}
}

func readEvalPrint(ivy *Ivy, l *liner.State, origMode, linerMode liner.ModeApplier) error {
	input, err := l.Prompt("ivy> ")
	if err != nil {
		return err
	}
	if specialDemo.MatchString(input) {
		err = origMode.ApplyMode()
		if err != nil {
			return err
		}
		runDemo()
		err = linerMode.ApplyMode()
		if err != nil {
			return err
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
		_, err = fmt.Fprintln(ivy.in, input)
		if err != nil {
			return err
		}
		fmt.Print(ivy.readResponse())
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
		fmt.Println("Could not run demo:", err)
	}
}

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
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
