/*
 * Load stringfilters from config file
 */

package main

import (
	"embed"
	"github.com/kghose/radio-go-go/internal"
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
		return filters 
	}

	data, err := loadData(path)
	if len(data) == 0 {
		data, _ = f.ReadFile("filters.json")
		overwriteData(path, data)
	}

	return filters
}
