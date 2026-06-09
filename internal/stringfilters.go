/*
 * Filter out strings or bits of strings we don't want
 */
package radio

import (
	"strings"
)

type StringFilters struct {
	StationNameJunkChars string
	NotSongTitles        []string
}

func (sf *StringFilters) IsSongTitle(song string) bool {
	for _, prefix := range sf.NotSongTitles {
		if strings.HasPrefix(song, prefix) {
			return false
		}
	}
	return true
}
