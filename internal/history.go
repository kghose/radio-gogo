/*
Largely, manage a list of stations
*/
package radio

import (
	"encoding/json"
	"io/fs"
	"log/slog"
	"os"
	"path"
)

const (
	userDataDir  = "radio-gogo"
	stationsFile = "stations.json"
)


func LoadHistory() ([]Station, error) {
	stations := []Station{}

	fname, err := stationsFilePath()
	if err != nil {
		return stations, err
	}

	dataBytes, err := os.ReadFile(fname)
	if os.IsNotExist(err) {
		slog.Info("No history file", "path", fname)
		return stations, nil
	}
	if err != nil {
		slog.Info("Error loading history file", "Error", err)
		return stations, err
	}

	if err = json.Unmarshal(dataBytes, &stations); err != nil {
		slog.Error("Error parsing history file", "Error", err)
		return stations, err
	}

	return stations, nil
}

func SaveHistory(stations []Station) error {
	fname, err := stationsFilePath()
	if err != nil {
		return err
	}

	f, err := os.Create(fname)
	if err != nil {
		slog.Error("Couldn't create history file", "Error", err)
		return err
	}
	defer f.Close()

	dataStr, err := json.MarshalIndent(stations, "", " ")
	if err != nil {
		slog.Error("Error formatting station history", "Error", err)
		return err
	}
	_, err = f.Write(dataStr)
	if err != nil {
		slog.Error("Error saving history", "Error", err)
		return err
	}

	return nil
}



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

