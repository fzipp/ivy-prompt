// Copyright 2022 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/peterh/liner"
)

func loadHistory(l *liner.State) error {
	path, err := historyPath()
	if err != nil {
		return fmt.Errorf("no config directory to load history from: %w", err)
	}
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open history file: %w", err)
	}
	defer f.Close()
	_, err = l.ReadHistory(f)
	if err != nil {
		return fmt.Errorf("could not read history file: %w", err)
	}
	return nil
}

func saveHistory(l *liner.State) error {
	path, err := historyPath()
	if err != nil {
		return fmt.Errorf("no config directory to save history: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create history file: %w", err)
	}
	defer f.Close()
	_, err = l.WriteHistory(f)
	if err != nil {
		return fmt.Errorf("could not write history file: %w", err)
	}
	return nil
}

func historyPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("could not determine user config directory: %w", err)
	}
	ivyDir := filepath.Join(configDir, "ivy")
	err = os.MkdirAll(ivyDir, 0o700)
	if err != nil {
		return "", fmt.Errorf("could not create ivy configuration directory: %w", err)
	}
	return filepath.Join(ivyDir, "ivy_history"), nil
}
