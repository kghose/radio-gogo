/*
Largely, manage loading and saving the station history to file.
*/
package main

import (
	"encoding/json"
	"log/slog"

	radio "github.com/kghose/radio-go-go/internal"
)

var historyPathConfig = ConfigPath{
	env:      "XDG_DATA_HOME",
	fallback: []string{".local", "share"},
	name:     "stations.json",
}

func LoadHistory() (map[string]*radio.Station, error) {
	stations := make(map[string]*radio.Station)

	path, err := getPath(historyPathConfig)
	if err != nil {
		return stations, err
	}

	data, err := loadData(path)
	if err != nil {
		return stations, err
	}

	if err = json.Unmarshal(data, &stations); err != nil {
		slog.Error("Error parsing history file", "Error", err)
		return stations, err
	}

	return stations, nil
}

func SaveHistory(stations map[string]*radio.Station) error {
	path, err := getPath(historyPathConfig)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(radio.History(stations), "", " ")
	if err != nil {
		slog.Error("Error formatting station history", "Error", err)
		return err
	}

	return overwriteData(path, data)
}
