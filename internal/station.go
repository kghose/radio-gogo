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
	Details    radiobrowser.Station
	LastPlayed time.Time // The last time we played this
	Favorite   bool      // In our favorites list?
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
	s.Details.Name = strings.TrimLeft(sanitize(&s.Details.Name), STATION_NAME_JUNK_CHARS)
	s.Details.URLResolved = sanitize(&s.Details.URLResolved)
	s.Details.URL = sanitize(&s.Details.URL)
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

func Search(index map[string]*Station, urls []string) map[string]*Station {
	s := make(map[string]*Station)
	for _, url := range urls {
		s[url] = index[url]
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
			strings.ToLower(a.Details.Name),
			strings.ToLower(b.Details.Name))
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

func MakeNewIndexFromSearch(
	sl []radiobrowser.Station,
	oldIndex map[string]*Station,
) (map[string]*Station, []string) {
	index := History(oldIndex)
	urls := []string{}
	for i := range sl {
		url := sl[i].URLResolved
		if _, ok := index[url]; ok {
			// TODO: Take care of station metadata changes
		} else {
			index[url] = &Station{sl[i], time.Time{}, false}
			sanitizeStation(index[url])
		}
		urls = append(urls, url)
	}
	return index, urls
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
