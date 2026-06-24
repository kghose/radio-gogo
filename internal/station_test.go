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
	sA := Station{
		Station: radiobrowser.Station{
			Name:        "a",
			URLResolved: "urlA",
			URL:         "boo",
		},
		LastPlayed: time.Date(2026, time.May, 12, 0, 0, 0, 0, time.UTC),
	}
	sB := Station{
		Station: radiobrowser.Station{
			Name:        "b",
			URLResolved: "urlB",
			URL:         "boo",
		},
		Favorite: true}
	sC := Station{
		Station: radiobrowser.Station{
			Name:        "c",
			URLResolved: "urlC",
			URL:         "boo",
		},
	}
	got := History(map[string]*Station{
		"urlA": &sA, "urlB": &sB, "urlC": &sC})
	want := map[string]*Station{
		"urlA": &sA, "urlB": &sB}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("History from index mismatch (-want +got):\n%s", diff)
	}
}

func TestMakeNewIndexFromSearch(t *testing.T) {
	iA := Station{
		Station: radiobrowser.Station{
			Name:        "a",
			URLResolved: "urlA",
			URL:         "boo",
		},
		LastPlayed: time.Date(2026, time.May, 12, 0, 0, 0, 0, time.UTC),
	}
	iB := Station{
		Station: radiobrowser.Station{
			Name:        "b",
			URLResolved: "urlB",
			URL:         "boo",
		},
		Favorite: true}
	iC := Station{
		Station: radiobrowser.Station{
			Name:        "c",
			URLResolved: "urlC",
			URL:         "boo",
		},
	}

	sN := radiobrowser.Station{
		Name:        "a-new",
		URLResolved: "urlA",
		URL:         "boo",
	}
	sM := radiobrowser.Station{
		Name:        "b",
		URLResolved: "urlM",
		URL:         "boo",
	}

	gotIndex, gotUrl := MakeNewIndexFromSearch(
		[]radiobrowser.Station{sN, sM},
		map[string]*Station{
			"urlA": &iA, "urlB": &iB, "urlC": &iC})

	iM := Station{Station: sM, SearchResult: true}
	wantIndex := map[string]*Station{
		"urlA": &iA, "urlB": &iB, "urlM": &iM}
	if diff := cmp.Diff(wantIndex, gotIndex); diff != "" {
		t.Errorf("New index from search mismatch (-want +got):\n%s", diff)
	}

	wantUrl := []string{"urlA", "urlM"}
	if diff := cmp.Diff(wantUrl, gotUrl); diff != "" {
		t.Errorf("New index from search URL mismatch (-want +got):\n%s", diff)
	}
}
