/*
Where to save our program data and logs
*/

package main

import (
	"io/fs"
	"log/slog"
	"os"
	"path"
)

const (
	programDir = "radio-gogo"
)

type ConfigPath struct {
	env      string
	fallback []string
	name     string
}

func getPath(cp ConfigPath) (string, error) {
	dir := os.Getenv(cp.env)
	if dir == "" {
		home := os.Getenv("HOME")
		dir = path.Join(home, path.Join(cp.fallback...), programDir)
	} else {
		dir = path.Join(dir, programDir)
	}
	path := path.Join(dir, cp.name)
	_, err := os.Stat(path)
	if err != nil {
		err = os.MkdirAll(dir, fs.FileMode(0777))
		if err != nil {
			slog.Error("Unable to create dir", "path", dir, "error", err)
			return "", err
		}
	}
	return path, err
}

func loadData(path string) ([]byte, error) {
	bytes, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		slog.Info("No file", "path", path)
	}
	if err != nil {
		slog.Info("Error loading file", "error", err)
	}
	return bytes, err
}

func overwriteData(path string, bytes []byte) error {
	f, err := os.Create(path)
	if err != nil {
		slog.Error("Unable to create file", "path", path, "error", err)
		return err
	}
	defer f.Close()

	_, err = f.Write(bytes)
	if err != nil {
		slog.Error("Error saving default data", "path", path, "error", err)
		return err
	}

	return err
}

func getHistoryPath() (string, error) {
	return getPath(historyPathConfig)
}
