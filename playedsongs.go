/*
 * Append the song buffer to a text file
 */

package main

import (
	"bufio"
	"log/slog"
	"os"
	"slices"

	"github.com/kghose/radio-go-go/internal"
)

type PlayedSongs struct {
	hist    radio.HistorySet
	filters *radio.StringFilters
	path    string
}

func (ps *PlayedSongs) load() {
	if ps.path == "" {
		return
	}

	fp, err := os.Open(ps.path)
	if err != nil {
		slog.Info("Could not open played songs file", "path", ps.path)
		return
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		ps.hist.Add(scanner.Text())
	}
	slog.Info("Loaded song title history", "N", ps.hist.Len())
}

func (ps *PlayedSongs) save() {
	if ps.path == "" {
		return
	}

	fp, err := os.Create(ps.path)
	if err != nil {
		slog.Error("Could not save layed songs file", "path", ps.path)
		ps.path = "" // stop trying
		return
	}
	defer fp.Close()

	writer := bufio.NewWriter(fp)
	for s := range ps.hist.ByOldest() {
		writer.WriteString(s + "\n")
	}
	err = writer.Flush()
	if err != nil {
		slog.Error("Problem saving songs", "err", err)
	}
}

func NewPlayedSongs(n int, filters *radio.StringFilters) PlayedSongs {

	ps := PlayedSongs{}
	ps.hist = radio.NewHistorySet(n)
	ps.filters = filters

	var playedSongsPathConfig = ConfigPath{
		env:      "XDG_DATA_HOME",
		fallback: []string{".local", "share"},
		name:     "songs.txt",
	}

	path, err := getPath(playedSongsPathConfig)
	if err != nil {
		slog.Error("Could not setup path for saving played song titles")
	} else {
		ps.path = path
		slog.Info("Saving played song titles", "path", ps.path)
	}
	ps.load()

	return ps
}

func (ps *PlayedSongs) add(song string) bool {
	if song != "" && ps.filters.IsSongTitle(song) && ps.hist.Latest() != song {
		ps.hist.Add(song)
		return true
	}
	return false
}

func (ps *PlayedSongs) latestPlayed() []string {
	return slices.Collect(ps.hist.ByLatest())
}
