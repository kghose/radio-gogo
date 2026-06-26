/*
 * Load stringfilters from config file
 */

package main

import (
	"embed"
	"encoding/json"
	"github.com/kghose/radio-go-go/internal"
	"log/slog"
)

//go:embed filters.json
var f embed.FS

func loadStringFilters() radio.StringFilters {

	var filtersPathConfig = ConfigPath{
		env:      "XDG_CONFIG_HOME",
		fallback: []string{".config"},
		name:     "filters.json",
	}

	filters := radio.StringFilters{}

	path, err := getPath(filtersPathConfig)
	if err != nil {
		slog.Error("Error loading filters file", "Error", err)
		return filters
	}

	data, err := loadData(path)
	if len(data) == 0 {
		data, _ = f.ReadFile("filters.json")
		slog.Info("Creating default string filters file")
		overwriteData(path, data)
	}

	if err = json.Unmarshal(data, &filters); err != nil {
		slog.Error("Error parsing string filters file", "Error", err)
	}

	return filters
}
