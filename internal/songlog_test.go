package radio

import (
	"fmt"
	"testing"

	cmp "github.com/google/go-cmp/cmp"
)

func TestAddFirstSong(t *testing.T) {
	sl := SongLog{}
	if sl.Add("New song") != true {
		t.Errorf("Adding new song doesn't return true")
	}
}

func TestAddSongInMiddle(t *testing.T) {
	sl := SongLog{}
	sl.Add("New song")
	if sl.Add("Another new song") != true {
		t.Errorf("Adding second new song doesn't return true")
	}
}

func TestAddSongWraparound(t *testing.T) {
	sl := SongLog{}
	for i := range songlogBufsize {
		song := fmt.Sprintf("Song %d", i)
		sl.Add(song)
	}
	if sl.Add("Another new song") != true {
		t.Errorf("Adding second new song doesn't return true")
	}
}

func TestAddDuplicateSong(t *testing.T) {
	sl := SongLog{}
	sl.Add("Same song")
	if sl.Add("Same song") != false {
		t.Errorf("Adding duplicate song doesn't return false")
	}
}

func TestAddDuplicateSongWraparound(t *testing.T) {
	sl := SongLog{}
	for i := range songlogBufsize {
		song := fmt.Sprintf("Song %d", i)
		sl.Add(song)
	}
	if sl.Add(fmt.Sprintf("Song %d", songlogBufsize - 1)) != false {
		t.Errorf("Adding duplicate song at end of buffer doesn't return false")
	}
}

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

			want := make([]string, songlogBufsize+tt.excess)
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
