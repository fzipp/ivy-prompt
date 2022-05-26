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
		return nil, err
	}
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	cmd.Stderr = cmd.Stdout
	out := bufio.NewReader(outPipe)
	return &Ivy{cmd: cmd, in: in, out: out, internalPrompt: '\x17'}, nil
}

func (ivy *Ivy) Start() error {
	err := ivy.cmd.Start()
	if err != nil {
		return err
	}
	return ivy.setInternalPrompt()
}

func (ivy *Ivy) Exec(stmt string) (string, error) {
	_, err := fmt.Fprintln(ivy.in, stmt)
	if err != nil {
		return "", err
	}
	return ivy.readResponse(), nil
}

func (ivy *Ivy) Ops() ([]string, error) {
	output, err := ivy.Exec(")op")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(output, "\n")
	var ops []string
	for _, line := range lines {
		if line != "" && line[0] == '\t' {
			ops = append(ops, line[1:])
		}
	}
	return ops, nil
}

func (ivy *Ivy) Quit() error {
	err := ivy.cmd.Process.Signal(os.Kill)
	if err != nil {
		return err
	}
	return ivy.cmd.Wait()
}

func (ivy *Ivy) setInternalPrompt() error {
	_, err := ivy.Exec(fmt.Sprintf(`)prompt "%c"`, ivy.internalPrompt))
	return err
}

func (ivy *Ivy) readResponse() string {
	s, _ := ivy.out.ReadString(ivy.internalPrompt)
	return s[:len(s)-2]
}
