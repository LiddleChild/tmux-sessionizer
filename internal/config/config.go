// Package config
package config

import (
	"bufio"
	"errors"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

const (
	AppName    = "tmux-sessionizer"
	AppVersion = "v1.0.2"
)

var (
	BaseConfigPath = path.Join(os.Getenv("HOME"), ".config", AppName)
	WorkspacesPath = path.Join(BaseConfigPath, "workspaces")
	EntriesPath    = path.Join(BaseConfigPath, "entries")

	WorkspaceEntries = make([]WorkspaceEntry, 0)
)

func Init() error {
	_, err := os.Stat(BaseConfigPath)
	switch {
	case err != nil && errors.Is(err, os.ErrNotExist):
		if err := os.MkdirAll(BaseConfigPath, 0o755); err != nil {
			return err
		}

	case err != nil:
		return err
	}

	workspaceEntries, err := readWorkspaceEntries()
	if err != nil {
		return err
	}

	singleEntries, err := readSingleEntries()
	if err != nil {
		return err
	}

	entriesSeq := slices.Values(append(workspaceEntries, singleEntries...))

	WorkspaceEntries = slices.SortedStableFunc(entriesSeq, func(a, b WorkspaceEntry) int {
		return strings.Compare(strings.ToLower(a.Path), strings.ToLower(b.Path))
	})

	return nil
}

func readWorkspaceEntries() ([]WorkspaceEntry, error) {
	workspaceEntries := []WorkspaceEntry{}

	f, err := os.Open(WorkspacesPath)
	switch {
	case err != nil && errors.Is(err, os.ErrNotExist):
		return workspaceEntries, nil
	case err != nil:
		return workspaceEntries, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		workspaceBasePath := absolutePath(scanner.Text())
		entries, err := os.ReadDir(workspaceBasePath)
		if err != nil {
			return workspaceEntries, err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				workspaceEntries = append(workspaceEntries, WorkspaceEntry{
					Name: strings.ReplaceAll(entry.Name(), ".", "_"),
					Path: path.Join(workspaceBasePath, entry.Name()),
				})
			}
		}
	}

	return workspaceEntries, nil
}

func readSingleEntries() ([]WorkspaceEntry, error) {
	singleEntries := []WorkspaceEntry{}

	f, err := os.Open(EntriesPath)
	switch {
	case err != nil && errors.Is(err, os.ErrNotExist):
		return singleEntries, nil
	case err != nil:
		return singleEntries, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		entryPath := absolutePath(scanner.Text())
		singleEntries = append(singleEntries, WorkspaceEntry{
			Name: strings.ReplaceAll(filepath.Base(entryPath), ".", "_"),
			Path: entryPath,
		})
	}

	return singleEntries, nil
}

func absolutePath(relative string) string {
	if strings.HasPrefix(relative, "~") {
		return path.Join(os.Getenv("HOME"), relative[2:]) // [2:] to remove ~/
	}

	return relative
}
