/*
Manage lists of stations
*/
package radio

import (
	radioBrowser "github.com/kghose/radio-go-go/internal/radio_browser"
	"slices"
	"sort"
	"time"
)

// A radio station we may have played and may have marked as favorite.
type Station struct {
	Details    radioBrowser.Station
	LastPlayed time.Time // The last time we played this
	Favorite   bool      // In our favorites list?
}

// Given a list of stations retrieved from a Radio Browser search and one loaded from
// our listening history, tag stations by last played time, and if they are favorites
// and return the list sorted by favorites, last played and then alphabetically
func SearchResults(
	stationsFromBrowser []radioBrowser.Station,
	stationsFromHistory []Station) []Station {
	historyMap := make(map[string]Station)
	for _, station := range stationsFromHistory {
		historyMap[station.Details.URLResolved] = station
	}
	searchResult := []Station{}
	for _, station := range stationsFromBrowser {
		if sta, ok := historyMap[station.URLResolved]; ok {
			// TODO: Take care of station metadata changes
			searchResult = append(searchResult, sta)
		} else {
			searchResult = append(searchResult, Station{station, time.Time{}, false})
		}
	}
	sortList(searchResult)
	return searchResult
}

func AddToHistory(station Station, history []Station) []Station {
	history = slices.Insert(history, 0, station)
	sortList(history)
	return history
}

func sortList(stations []Station) {
	sort.SliceStable(stations, func(i, j int) bool {
		if stations[i].Favorite == stations[j].Favorite {
			cmp := stations[i].LastPlayed.Compare(stations[j].LastPlayed)
			if cmp == 0 {
				return stations[i].Details.Name < stations[j].Details.Name
			}
			return cmp > 0
		} else {
			return stations[i].Favorite
		}
	})
}
