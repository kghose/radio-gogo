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

// The url we just played can come from the history list or the search list
// Even if it comes from the history list, if we find it in the search results
// we update the history, assuming that the search results have the most upto date
// metadata about the station
func AddToHistory(url string, searchResult []Station, history []Station) []Station {
	stationDetails := radioBrowser.Station{}
	for i := range searchResult {
		if searchResult[i].Details.URLResolved == url {
			stationDetails = searchResult[i].Details
			break
		}
	}
	for i := range history {
		if history[i].Details.URLResolved == url {
			if stationDetails.URLResolved == url {
				// The station is in the search and history
				// update the details and the last played time
				history[i].Details = stationDetails
				history[i].LastPlayed = time.Now()
				return history
			} else {
				// It's just in the history
				history[i].LastPlayed = time.Now()
				return history
			}
		}
	}

	// The station is not in the history
	newStation := Station{stationDetails, time.Now(), false}
	history = slices.Insert(history, 0, newStation)
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
