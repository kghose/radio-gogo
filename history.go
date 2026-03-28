/*
Largely, manage a list of stations
*/
package main 

import (
	"encoding/json"
	"log/slog"
	"os"

	radio "github.com/kghose/radio-go-go/internal"
)

func LoadHistory() ([]radio.Station, error) {
	stations := []radio.Station{}

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

func SaveHistory(stations []radio.Station) error {
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
