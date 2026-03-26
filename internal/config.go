/*
Where to save our program data and logs
*/

package radio

import (
	"io/fs"
	"log/slog"
	"os"
	"path"
)

const (
	userDataDir  = "radio-gogo"
	stationsFile = "stations.json"
	logFile = "log.txt"
)

func stationsFilePath() (string, error) {
	home := os.Getenv("HOME")
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		dataHome = path.Join(home, ".local", "share")
	}
	dataHome = path.Join(dataHome, userDataDir)
	err := os.MkdirAll(dataHome, fs.FileMode(0777))
	if err != nil {
		slog.Error("Unable to create data directory", "Error", err)
	}
	return path.Join(
		dataHome,
		stationsFile,
	), err
}

func logsFilePath() (string, error) {
	home := os.Getenv("HOME")
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		dataHome = path.Join(home, ".local", "share")
	}
	dataHome = path.Join(dataHome, userDataDir)
	err := os.MkdirAll(dataHome, fs.FileMode(0777))
	if err != nil {
		slog.Error("Unable to create data directory", "Error", err)
	}
	return path.Join(
		dataHome,
		logFile,
	), err
}
