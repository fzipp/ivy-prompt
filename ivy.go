// Copyright 2022 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Ivy struct {
	cmd *exec.Cmd
	in  io.Writer
	out *bufio.Reader

	internalPrompt byte
}

func NewIvy() (*Ivy, error) {
	cmd := exec.Command("ivy")
	in, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("could not connect to ivy input: %w", err)
	}
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("could not connect to ivy output: %w", err)
	}
	cmd.Stderr = cmd.Stdout
	out := bufio.NewReader(outPipe)
	return &Ivy{cmd: cmd, in: in, out: out, internalPrompt: '\x17'}, nil
}

func (ivy *Ivy) Start() error {
	err := ivy.cmd.Start()
	if err != nil {
		return fmt.Errorf("could not start ivy: %w", err)
	}
	return ivy.setInternalPrompt()
}

func (ivy *Ivy) Exec(stmt string) (string, error) {
	_, err := fmt.Fprintln(ivy.in, stmt)
	if err != nil {
		return "", fmt.Errorf("could not send input to ivy: %w", err)
	}
	return ivy.readResponse()
}

func (ivy *Ivy) Ops() ([]string, error) {
	return ivy.listDefs("user-defined ops", ")op")
}

func (ivy *Ivy) Vars() ([]string, error) {
	return ivy.listDefs("user-defined vars", ")var")
}

func (ivy *Ivy) listDefs(description, command string) ([]string, error) {
	output, err := ivy.Exec(command)
	if err != nil {
		return nil, fmt.Errorf("could not get list of %s: %w", description, err)
	}
	lines := strings.Split(output, "\n")
	var defs []string
	for _, line := range lines {
		if line != "" && line[0] == '\t' {
			defs = append(defs, line[1:])
		}
	}
	return defs, nil
}

func (ivy *Ivy) Quit() error {
	err := ivy.cmd.Process.Signal(os.Kill)
	if err != nil {
		return fmt.Errorf("could not quit ivy process: %w", err)
	}
	return ivy.cmd.Wait()
}

func (ivy *Ivy) setInternalPrompt() error {
	_, err := ivy.Exec(fmt.Sprintf(`)prompt "%c"`, ivy.internalPrompt))
	if err != nil {
		return fmt.Errorf("could not set internal prompt: %w", err)
	}
	return nil
}

func (ivy *Ivy) readResponse() (string, error) {
	s, err := ivy.out.ReadString(ivy.internalPrompt)
	if err != nil {
		return "", fmt.Errorf("could not read response from ivy: %w", err)
	}
	return s[:len(s)-2], nil
}
