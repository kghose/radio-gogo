package radio

import (
	cmp "github.com/google/go-cmp/cmp"
	radioBrowser "github.com/kghose/radio-go-go/internal/radio_browser"
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
		radioBrowser.Station{
			Name:        "\x00name",
			URLResolved: "url\x0d",
			URL:         "\u200burl",
		},
		time.Time{},
		true}
	sanitizeStation(&station)
	want := Station{
		radioBrowser.Station{
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
		radioBrowser.Station{
			Name:        "a",
			URLResolved: "urlA",
			URL:         "boo",
		},
		time.Date(2026, time.May, 12, 0, 0, 0, 0, time.UTC),
		false}
	sB := Station{
		radioBrowser.Station{
			Name:        "b",
			URLResolved: "urlB",
			URL:         "boo",
		},
		time.Time{},
		true}
	sC := Station{
		radioBrowser.Station{
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
		radioBrowser.Station{
			Name:        "a",
			URLResolved: "urlA",
			URL:         "boo",
		},
		time.Date(2026, time.May, 12, 0, 0, 0, 0, time.UTC),
		false}
	iB := Station{
		radioBrowser.Station{
			Name:        "b",
			URLResolved: "urlB",
			URL:         "boo",
		},
		time.Time{},
		true}
	iC := Station{
		radioBrowser.Station{
			Name:        "c",
			URLResolved: "urlC",
			URL:         "boo",
		},
		time.Time{},
		false}

	sN := radioBrowser.Station{
		Name:        "a-new",
		URLResolved: "urlA",
		URL:         "boo",
	}
	sM := radioBrowser.Station{
		Name:        "b",
		URLResolved: "urlM",
		URL:         "boo",
	}

	gotIndex, gotUrl := MakeNewIndexFromSearch(
		[]radioBrowser.Station{sN, sM},
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
