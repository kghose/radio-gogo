/*
Manage collections of stations.

We create one set of Station objects and share them across the multiple lists we
manage.
*/
package radio

import (
	"github.com/kghose/radio-go-go/internal/radiobrowser"

	"slices"
	"strings"
	"time"
	"unicode"
)

const STATION_NAME_JUNK_CHARS = "_.-+*# "

// A radio station we may have played and may have marked as favorite.
type Station struct {
	radiobrowser.Station
	LastPlayed   time.Time // The last time we played this
	Favorite     bool      // In our favorites list?
	SearchResult bool      `json:"-"` // Current search result? (No need to save in station history)
}

func sanitize(s *string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, *s)
}

func sanitizeStation(s *Station) {
	s.Name = strings.TrimLeft(sanitize(&s.Name), STATION_NAME_JUNK_CHARS)
	s.URLResolved = sanitize(&s.URLResolved)
	s.URL = sanitize(&s.URL)
}

func History(index map[string]*Station) map[string]*Station {
	history := make(map[string]*Station)
	for k, v := range index {
		if !v.LastPlayed.IsZero() || v.Favorite {
			history[k] = v
		}
	}
	return history
}

func Faves(index map[string]*Station) map[string]*Station {
	faves := make(map[string]*Station)
	for k, v := range index {
		if v.Favorite {
			faves[k] = v
		}
	}
	return faves
}

func Search(index map[string]*Station) map[string]*Station {
	s := make(map[string]*Station)
	for k, v := range index {
		if v.SearchResult {
			s[k] = v
		}
	}
	return s
}

func SortAlpha(index map[string]*Station) []*Station {
	l := []*Station{}
	for _, v := range index {
		l = append(l, v)
	}
	slices.SortFunc(l, func(a, b *Station) int {
		return strings.Compare(
			strings.ToLower(a.Name),
			strings.ToLower(b.Name))
	})
	return l
}

func SortLastPlayed(index map[string]*Station) []*Station {
	l := []*Station{}
	for _, v := range index {
		l = append(l, v)
	}
	slices.SortFunc(l, func(a, b *Station) int {
		return b.LastPlayed.Compare(a.LastPlayed)
	})
	return l
}

func wipeSearchFlagFromIndex(index map[string]*Station) {
	for _, v := range index {
		v.SearchResult = false
	}
}

func MakeNewIndexFromSearch(
	sl []radiobrowser.Station,
	oldIndex map[string]*Station,
) map[string]*Station {
	index := History(oldIndex)
	wipeSearchFlagFromIndex(index)
	for i := range sl {
		url := sl[i].URLResolved
		if _, ok := index[url]; ok {
			// TODO: Take care of station metadata changes
			index[url].SearchResult = true
		} else {
			index[url] = &Station{
				Station:      sl[i],
				SearchResult: true,
			}
			sanitizeStation(index[url])
		}
	}
	return index
}

type StationOp string

const (
	PLAYED StationOp = "Played"
	FAVE   StationOp = "Faved"
	UNFAVE StationOp = "Unfaved"
)

func UpdateIndex(url string, index map[string]*Station, op StationOp) {
	switch op {
	case PLAYED:
		index[url].LastPlayed = time.Now()
	case FAVE:
		index[url].Favorite = true
	case UNFAVE:
		index[url].Favorite = false
	}
}
