package radio

import (
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAddFirstSong(t *testing.T) {
	sl := NewSongLog(3)
	if sl.Add("New song") != true {
		t.Errorf("Adding new song doesn't return true")
	}
}

func TestAddSongInMiddle(t *testing.T) {
	sl := NewSongLog(3)
	sl.Add("New song")
	if sl.Add("Another new song") != true {
		t.Errorf("Adding second new song doesn't return true")
	}
}

func TestAddSongWraparound(t *testing.T) {
	sl := NewSongLog(2)
	sl.Add("Song 1")
	sl.Add("Song 2")
	if sl.Add("Song 1") != true {
		t.Errorf("Adding new song on wrap around doesn't return true")
	}
}

func TestAddDuplicateSong(t *testing.T) {
	sl := NewSongLog(2)
	sl.Add("Song 1")
	if sl.Add("Song 1") != false {
		t.Errorf("Adding duplicate song doesn't return false")
	}
}

func TestAddDuplicateSongWraparound(t *testing.T) {
	sl := NewSongLog(2)
	sl.Add("Song 1")
	sl.Add("Song 2")
	if sl.Add("Song 2") != false {
		t.Errorf("Adding duplicate song on wraparound doesn't return false")
	}
}

func TestRetrieval(t *testing.T) {
	tests := []struct {
		name  string
		songs []string
	}{
		{"No wrap around", []string{"a", "b"}},
		{"Single wrap around", []string{"a", "b", "c"}},
		{"Multiple wrap around", []string{"a", "b", "c", "d", "e"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			sl := NewSongLog(2)

			want := slices.Clone(tt.songs)
			slices.Reverse(want)
			want = want[:2]
			for _, s := range tt.songs {
				sl.Add(s)
			}

			got := sl.Songs()

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf(
					"Song log mismatch (-want +got):\n%s",
					diff)
			}
		})
	}
}
