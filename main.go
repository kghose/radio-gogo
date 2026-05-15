package main

import (
	"log/slog"
	"time"

	radio "github.com/kghose/radio-go-go/internal"
	mpv "github.com/kghose/radio-go-go/internal/mpv"
	radio_browser "github.com/kghose/radio-go-go/internal/radio_browser"
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

	shs := NewSongHistorySaver()

	stationIndex, err := LoadHistory()
	// TODO: Handle errors
	if err != nil {
	}
	searchUrls := []string{}
	keywords := "Search"

	ui := radio.UI{}

	stnFunc := func(op radio.StationOp) {
		radio.UpdateIndex(ui.SelectedURL(), stationIndex, op)
		ui.RefreshLists(stationIndex, searchUrls, keywords)
		SaveHistory(stationIndex)
	}

	keyMap := map[rune]radio.KeyFunc{
		'h': {Help: "Show history pane", Fn: ui.ShowHist},
		's': {Help: "Show search pane", Fn: ui.ShowSearch},
		'/': {Help: "Search", Fn: ui.ShowSearchBar},
		'f': {Help: "Show faves pane", Fn: ui.ShowFaves},
		'=': {Help: "Fave station", Fn: func() { stnFunc(radio.FAVE) }},
		'-': {Help: "Unfave station", Fn: func() { stnFunc(radio.UNFAVE) }},
		'p': {Help: "Pause", Fn: func() { mpvPlayer.TogglePause() }},
		'q': {Help: "Quit", Fn: ui.Stop},
	}

	searchFunc := func(kw string) {
		keywords = kw
		stations, err := radio_browser.StationSearch(keywords, server)
		if err != nil {
			slog.Error("Error searching for stations.")
		}
		stationIndex, searchUrls =
			radio.MakeNewIndexFromSearch(stations, stationIndex)
		ui.RefreshLists(stationIndex, searchUrls, keywords)
		ui.ResetSearchScroll()
		ui.ShowSearch()
	}

	playFunc := func(idx int, _ string, url string, _ rune) {
		r := mpvPlayer.Play(url)
		slog.Info("Play", "url", url, "mpv", r.Error)
		stnFunc(radio.PLAYED)
	}

	periodicInfoRefreshFunc := func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for ; ; <-ticker.C {
			meta := mpvPlayer.Meta()
			ui.SetNowPlaying(meta)
			shs.save(meta.Title)
		}
	}

	ui.Setup(
		keyMap,
		searchFunc,
		playFunc,
	)

	go periodicInfoRefreshFunc()

	ui.RefreshLists(stationIndex, searchUrls, keywords)
	ui.ShowHist()
	if err := ui.Run(); err != nil {
		panic(err)
	}
}
