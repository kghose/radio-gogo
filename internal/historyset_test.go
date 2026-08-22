package radio

import (
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAddToHistorySet(t *testing.T) {
	tests := []struct {
		name   string
		inputs []string
		want   []string
	}{
		{
			"Unique<max",
			[]string{"a", "b", "c"},
			[]string{"c", "b", "a"},
		},
		{
			"Unique>max",
			[]string{"a", "b", "c", "d"},
			[]string{"d", "c", "b"},
		},
		{
			"Repeated<max",
			[]string{"a", "b", "a", "b"},
			[]string{"b", "a"},
		},
		{
			"Repeated>max",
			[]string{"a", "b", "c", "a"},
			[]string{"a", "c", "b"},
		},
		{
			"Repeated>>max",
			[]string{"a", "b", "c", "d", "a"},
			[]string{"a", "d", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hs := NewHistorySet(3)
			for _, inp := range tt.inputs {
				hs.Add(inp)
			}

			got := slices.Collect(hs.ByLatest())
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf(
					"History set mismatch (-want +got):\n%s",
					diff)
			}
		})
	}
}

func TestByOldest(t *testing.T) {
	hs := NewHistorySet(3)
	hs.Add("a")
	hs.Add("b")
	hs.Add("c")
	want := []string{"a", "b", "c"}

	got := slices.Collect(hs.ByOldest())
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf(
			"History set mismatch (-want +got):\n%s",
			diff)
	}

}

