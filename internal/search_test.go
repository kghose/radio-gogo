package radio

import (
	"testing"

	cmp "github.com/google/go-cmp/cmp"

	radioBrowser "github.com/kghose/radio-go-go/internal/radio_browser"
)

func TestParseSearchString(t *testing.T) {
	tests := []struct {
		name         string
		searchString string
		want         radioBrowser.SearchQuery
	}{
		{
			"No keys",
			"hard rock,jazz",
			radioBrowser.SearchQuery{
				TagList: []string{"hard rock,jazz"}}},
		{
			"Multiple keys",
			"t:hard rock t:jazz",
			radioBrowser.SearchQuery{
				TagList: []string{"hard rock", "jazz"}}},
		{
			"All keys",
			"n:BBC c:antigua c:Albania t:jazz  t: classic rock",
			radioBrowser.SearchQuery{
				Name:    "BBC",
				Country: "Albania",
				TagList: []string{"jazz", "classic rock"}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseSearchString(tt.searchString)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf(
					"Search string parse mismatch (-want +got):\n%s",
					diff)
			}
		})
	}
}
