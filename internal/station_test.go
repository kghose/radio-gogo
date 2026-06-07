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
		radiobrowser.Station{
			Name:        "\x00name",
			URLResolved: "url\x0d",
			URL:         "\u200burl",
		},
		time.Time{},
		true}
	sanitizeStation(&station)
	want := Station{
		radiobrowser.Station{
			Name:        "name",
			URLResolved: "url",
			URL:         "url",
		},
		time.Time{},
		true}
	if station != want {
		t.Errorf("Sanitize station fails")
	}
}

func TestHistory(t *testing.T) {
	sA := Station{
		radiobrowser.Station{
			Name:        "a",
			URLResolved: "urlA",
			URL:         "boo",
		},
		time.Date(2026, time.May, 12, 0, 0, 0, 0, time.UTC),
		false}
	sB := Station{
		radiobrowser.Station{
			Name:        "b",
			URLResolved: "urlB",
			URL:         "boo",
		},
		time.Time{},
		true}
	sC := Station{
		radiobrowser.Station{
			Name:        "c",
			URLResolved: "urlC",
			URL:         "boo",
		},
		time.Time{},
		false}
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
		radiobrowser.Station{
			Name:        "a",
			URLResolved: "urlA",
			URL:         "boo",
		},
		time.Date(2026, time.May, 12, 0, 0, 0, 0, time.UTC),
		false}
	iB := Station{
		radiobrowser.Station{
			Name:        "b",
			URLResolved: "urlB",
			URL:         "boo",
		},
		time.Time{},
		true}
	iC := Station{
		radiobrowser.Station{
			Name:        "c",
			URLResolved: "urlC",
			URL:         "boo",
		},
		time.Time{},
		false}

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

	iM := Station{sM, time.Time{}, false}
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
