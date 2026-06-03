package radio

import (
	"fmt"
	"testing"

	cmp "github.com/google/go-cmp/cmp"
)

func TestRetrieval(t *testing.T) {
	tests := []struct {
		name   string
		excess int
	}{
		{"No wrap around", 0},
		{"Single wrap around", songlogBufsize - 1},
		{"Multiple wrap around", 2*songlogBufsize - 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			sl := SongLog{}

			want := make([]string, songlogBufsize + tt.excess)
			for i := range songlogBufsize + tt.excess {
				song := fmt.Sprintf("Song %d", i)
				sl.Add(song)
				want[songlogBufsize+tt.excess-i-1] = song
			}

			got := []string{}
			for song := range sl.Songs() {
				got = append(got, song)
			}

			if diff := cmp.Diff(want[:songlogBufsize], got); diff != "" {
				t.Errorf(
					"Song log mismatch (-want +got):\n%s",
					diff)
			}
		})
	}
}
