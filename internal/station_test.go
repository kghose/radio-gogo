package radio

import (
	cmp "github.com/google/go-cmp/cmp"
	"github.com/kghose/radio-go-go/internal/radiobrowser"
	"testing"
	"time"
)

func TestSanitize(t *testing.T) {
	unprintable := "Put \x00another\x0a nickel\x0d in\n"
	got := sanitize(&unprintable)
	want := "Put another nickel in"
	if got != want {
		t.Errorf("Sanitize string fails")
	}
}

func TestSanitizeStation(t *testing.T) {
	station := Station{
		Station: radiobrowser.Station{
			Name:        "\x00name",
			URLResolved: "url\x0d",
			URL:         "\u200burl",
		},
		Favorite: true}
	sanitizeStation(&station)
	want := Station{
		Station: radiobrowser.Station{
			Name:        "name",
			URLResolved: "url",
			URL:         "url",
		},
		Favorite: true}
	if station != want {
		t.Errorf("Sanitize station fails")
	}
}

func TestHistory(t *testing.T) {
	playedStation := Station{
		Station: radiobrowser.Station{
			URLResolved: "played",
		},
		LastPlayed: time.Date(2026, time.May, 12, 0, 0, 0, 0, time.UTC),
	}
	favedStation := Station{
		Station: radiobrowser.Station{
			URLResolved: "faved",
		},
		Favorite: true}
	otherStation := Station{
		Station: radiobrowser.Station{
			URLResolved: "other",
		},
	}

	index := map[string]*Station{
		"played": &playedStation,
		"faved":  &favedStation,
		"other":  &otherStation,
	}
	want := map[string]*Station{
		"played": &playedStation,
		"faved":  &favedStation,
	}
	got := History(index)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("History from index mismatch (-want +got):\n%s", diff)
	}
}

func TestMakeNewIndexFromSearch(t *testing.T) {
	currentIndexA := Station{
		Station: radiobrowser.Station{
			URLResolved: "played",
		},
		LastPlayed: time.Date(2026, time.May, 12, 0, 0, 0, 0, time.UTC),
	}
	currentIndexB := Station{
		Station: radiobrowser.Station{
			URLResolved: "faved",
		},
		Favorite:     true,
		SearchResult: true,
	}
	currentIndexC := Station{
		Station: radiobrowser.Station{
			URLResolved: "search",
		},
		SearchResult: true,
	}
	currentIndex := map[string]*Station{
		"played": &currentIndexA,
		"faved":  &currentIndexB,
		"search": &currentIndexC,
	}

	sN := radiobrowser.Station{
		URLResolved: "played",
	}
	sM := radiobrowser.Station{
		URLResolved: "new search",
	}
	searchResult := []radiobrowser.Station{sN, sM}

	newIndexA := Station{
		Station: radiobrowser.Station{
			URLResolved: "played",
		},
		LastPlayed:   time.Date(2026, time.May, 12, 0, 0, 0, 0, time.UTC),
		SearchResult: true,
	}
	newIndexB := Station{
		Station: radiobrowser.Station{
			URLResolved: "faved",
		},
		Favorite: true,
	}
	newIndexC := Station{
		Station: radiobrowser.Station{
			URLResolved: "new search",
		},
		SearchResult: true,
	}
	want := map[string]*Station{
		"played":     &newIndexA,
		"faved":      &newIndexB,
		"new search": &newIndexC,
	}

	got := MakeNewIndexFromSearch(searchResult, currentIndex)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("New index from search mismatch (-want +got):\n%s", diff)
	}

}
