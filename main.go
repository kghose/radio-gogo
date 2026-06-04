package main

import (
	"log/slog"
	"time"

	radio "github.com/kghose/radio-go-go/internal"
	mpv "github.com/kghose/radio-go-go/internal/mpv"
	radio_browser "github.com/kghose/radio-go-go/internal/radio_browser"
)

type StationOp int

const (
	playStation StationOp = iota
	faveStation
	unfaveStation
	favePlayingStation
	unfavePlayingStation
)

func main() {
	slogger, closeFunc := SetupLoggingToFile()
	defer closeFunc()
	slog.SetDefault(slogger)

	mpvPlayer := mpv.Player{}
	mpvPlayer.Start()
	defer mpvPlayer.Quit()

	servers, err := radio_browser.GetAvailableServers()
	if err != nil {
		slog.Error("Could not find radio browser servers")
		servers = []string{""}
	}

	server := radio_browser.PickRandomServer(servers)

	songlog := radio.SongLog{}

	stationIndex, err := LoadHistory()
	// TODO: Handle errors
	if err != nil {
	}
	searchUrls := []string{}
	searchString := "Search"
	playingStationUrl := ""

	ui := radio.UI{}

	stnFunc := func(op StationOp) {
		url := ui.SelectedURL()
		var indexOp radio.StationOp
		switch op {
		case favePlayingStation:
			indexOp = radio.FAVE
			url = playingStationUrl
		case unfavePlayingStation:
			indexOp = radio.UNFAVE
			url = playingStationUrl
		case playStation:
			indexOp = radio.PLAYED
		case faveStation:
			indexOp = radio.FAVE
		case unfaveStation:
			indexOp = radio.UNFAVE
		}
		if url == "" {
			// Can only happen if we try to fave/unfave
			// current station and nothing has been played.
			return
		}

		radio.UpdateIndex(url, stationIndex, indexOp)
		ui.RefreshLists(stationIndex, searchUrls, searchString)
		SaveHistory(stationIndex)
	}

	keyMap := map[rune]radio.KeyFunc{
		'h': {Help: "Show history pane", Fn: ui.ShowHist},
		's': {Help: "Show search pane", Fn: ui.ShowSearch},
		'/': {Help: "Search", Fn: ui.ShowSearchBar},
		'f': {Help: "Show faves pane", Fn: ui.ShowFaves},
		'=': {Help: "Fave station", Fn: func() { stnFunc(faveStation) }},
		'-': {Help: "Unfave station", Fn: func() { stnFunc(unfaveStation) }},
		'+': {Help: "Fave playing station", Fn: func() { stnFunc(favePlayingStation) }},
		'_': {Help: "Unfave playing station", Fn: func() { stnFunc(unfavePlayingStation) }},
		'.': {Help: "Show played songs", Fn: ui.ShowPlayedsongs},
		'?': {Help: "Show help", Fn: ui.ShowHelp},
		'p': {Help: "Pause", Fn: func() { mpvPlayer.TogglePause() }},
		'q': {Help: "Quit", Fn: ui.Stop},
	}

	searchFunc := func(searchStr string) {
		searchString = searchStr
		sq := radio.ParseSearchString(searchStr)
		stations, err := radio_browser.StationSearch(sq, server)
		if err != nil {
			slog.Error("Error searching for stations.")
		}
		stationIndex, searchUrls =
			radio.MakeNewIndexFromSearch(stations, stationIndex)
		ui.RefreshLists(stationIndex, searchUrls, searchString)
		ui.ResetSearchScroll()
		ui.ShowSearch()
	}

	playFunc := func(idx int, _ string, url string, _ rune) {
		playingStationUrl = url
		r := mpvPlayer.Play(url)
		slog.Info("Play", "url", url, "mpv", r.Error)
		stnFunc(playStation)
	}

	periodicInfoRefreshFunc := func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for ; ; <-ticker.C {
			meta := mpvPlayer.Meta()
			ui.SetNowPlaying(meta)
			if songlog.Add(meta.Title) {
				ui.RefreshPlayedsongs(songlog.Songs())
			}
		}
	}

	ui.Setup(
		keyMap,
		searchFunc,
		playFunc,
	)

	go periodicInfoRefreshFunc()

	ui.RefreshLists(stationIndex, searchUrls, searchString)
	ui.ShowHist()
	if err := ui.Run(); err != nil {
		panic(err)
	}
}
