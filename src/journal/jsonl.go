// Copyright (C) 2026 Podomy.
// SPDX-License-Identifier: AGPL-3.0-or-later

package journal

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type JSONL struct {
	file    *os.File
	scanner *bufio.Scanner
	mu      sync.Mutex
}

// Append marshals the event and appends it as a JSON line to the journal file.
// It checks for context cancellation before acquiring the write lock.
func (j *JSONL) Append(ctx context.Context, event Event) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("append cancelled: %w", ctx.Err())
	default:
	}

	j.mu.Lock()
	defer j.mu.Unlock()

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal journal event: %w", err)
	}

	if _, err := j.file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("write journal event: %w", err)
	}

	if err := j.file.Sync(); err != nil {
		return fmt.Errorf("sync journal: %w", err)
	}

	return nil
}

// Close flushes and closes the underlying journal file.
func (j *JSONL) Close() error {
	if err := j.file.Close(); err != nil {
		return fmt.Errorf("close journal: %w", err)
	}

	return nil
}

// getJournalPath returns the auto-determined path for the local journal file.
func getJournalPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("get user config directory: %w", err)
	}

	appDir := filepath.Join(dir, "hive")
	if err := os.MkdirAll(appDir, 0o700); err != nil {
		return "", fmt.Errorf("create node config directory: %w", err)
	}

	return filepath.Join(appDir, "journal.jsonl"), nil
}

// OpenJSONL opens a JSONL journal file at the auto-determined path, creating it if it doesn't exist.
func OpenJSONL() (*JSONL, error) {
	path, err := getJournalPath()
	if err != nil {
		return nil, err
	}

	// #nosec G304: journal paths are local runtime configuration.
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open journal: %w", err)
	}

	// Creating a scanner for the file
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	jsonl := JSONL{
		file:    file,
		scanner: scanner,
	}

	return &jsonl, nil
}

func (j *JSONL) Read(ctx context.Context) (*Event, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("read cancelled: %w", ctx.Err())
	default:
	}

	j.mu.Lock()
	defer j.mu.Unlock()

	// Get a line from the file.
	if !j.scanner.Scan() {
		if err := j.scanner.Err(); err != nil {
			return nil, fmt.Errorf("scan journal: %w", j.scanner.Err())
		}
		return nil, io.EOF
	}
	line := j.scanner.Text()

	// Process the line.
	var event Event
	err := json.Unmarshal([]byte(line), &event)
	if err != nil {
		return nil, fmt.Errorf("unmarshal json: %w", err)
	}

	return &event, nil
}
