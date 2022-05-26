// Copyright 2022 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"path/filepath"

	"github.com/peterh/liner"
)

func loadHistory(l *liner.State) error {
	path, err := historyPath()
	if err != nil {
		return err
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = l.ReadHistory(f)
	if err != nil {
		return err
	}
	return nil
}

func saveHistory(l *liner.State) error {
	path, err := historyPath()
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = l.WriteHistory(f)
	if err != nil {
		return err
	}
	return nil
}

func historyPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	ivyDir := filepath.Join(configDir, "ivy")
	err = os.MkdirAll(ivyDir, 0700)
	if err != nil {
		return "", err
	}
	return filepath.Join(ivyDir, "ivy_history"), nil
}
