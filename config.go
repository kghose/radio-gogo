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
	userDataDirName = "radio-gogo"
	stationsFile    = "stations.json"
	songHistoryFile = "songs.txt"
	logFile         = "log.txt"
)

func getUserDataDirPath() (string, error) {
	home := os.Getenv("HOME")
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		dataHome = path.Join(home, ".local", "share")
	}
	userDataDirPath := path.Join(dataHome, userDataDirName)
	err := os.MkdirAll(userDataDirPath, fs.FileMode(0777))
	if err != nil {
		slog.Error("Unable to create data directory", "Error", err)
		return "", err
	}
	return userDataDirPath, nil
}

func stationsFilePath() (string, error) {
	userDataDirPath, err := getUserDataDirPath()
	return path.Join(
		userDataDirPath,
		stationsFile,
	), err
}

func songHistoryFilePath() (string, error) {
	userDataDirPath, err := getUserDataDirPath()
	return path.Join(
		userDataDirPath,
		songHistoryFile,
	), err
}

func logsFilePath() (string, error) {
	userDataDirPath, err := getUserDataDirPath()
	return path.Join(
		userDataDirPath,
		logFile,
	), err
}
