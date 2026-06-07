package radio

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/kghose/radio-go-go/internal/radiobrowser"
)

func TestParseSearchString(t *testing.T) {
	tests := []struct {
		name         string
		searchString string
		want         radiobrowser.SearchQuery
	}{
		{
			"No keys",
			"hard rock,jazz",
			radiobrowser.SearchQuery{
				TagList: []string{"hard rock,jazz"}}},
		{
			"Multiple keys",
			"t:hard rock t:jazz",
			radiobrowser.SearchQuery{
				TagList: []string{"hard rock", "jazz"}}},
		{
			"All keys",
			"n:BBC c:antigua c:Albania t:jazz  t: classic rock",
			radiobrowser.SearchQuery{
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
