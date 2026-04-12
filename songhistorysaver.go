package main

import (
	"log/slog"
	"os"
)

func songHistoryFp() (*os.File, error) {
	songHistoryFile, err := songHistoryFilePath()
	f, err := os.OpenFile(songHistoryFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		slog.Error("Couldn't open song history file for writing", "Error", err)
		return nil, err
	}
	return f, nil
}

type SongHistorySaver struct {
	fp       *os.File
	lastSong string
}

func NewSongHistorySaver() *SongHistorySaver {
	shs := SongHistorySaver{}
	shs.fp, _ = songHistoryFp()
	return &shs
}

func (shs *SongHistorySaver) save(song string) {
	if shs.fp == nil {
		return
	}
	if shs.lastSong != song {
		shs.lastSong = song
		if _, err := shs.fp.WriteString(shs.lastSong + "\n"); err != nil {
			slog.Error("Error writing song history", "error", err)
		}
		slog.Info("Playing", "song", song)
	}
}
