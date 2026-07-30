/*
 * Append the song buffer to a text file
 */

package main

import (
	"log/slog"
	"os"

	"github.com/kghose/radio-go-go/internal"
)

type PlayedSongs struct {
	fp      *os.File
	buf     radio.SongLog
	filters *radio.StringFilters
}

func NewPlayedSongs(n int, filters *radio.StringFilters) PlayedSongs {

	ps := PlayedSongs{}
	ps.buf = radio.NewSongLog(n)
	ps.filters = filters

	var playedSongsPathConfig = ConfigPath{
		env:      "XDG_DATA_HOME",
		fallback: []string{".local", "share"},
		name:     "songs.txt",
	}

	path, err := getPath(playedSongsPathConfig)
	if err != nil {
		slog.Error("Could not setup path for saving song files")
	} else {
		fp, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			slog.Error("Could not open played songs file")
		}
		ps.fp = fp
	}

	return ps
}

func (ps *PlayedSongs) add(song string) bool {
	if song != "" && ps.filters.IsSongTitle(song) && ps.buf.Add(song) {
		if ps.fp != nil {
			ps.fp.WriteString(song + "\n")
		}
		return true
	}
	return false
}
