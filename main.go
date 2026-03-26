package main

import (
	//	"fmt"
	//	"github.com/gdamore/tcell/v2"
	radio "github.com/kghose/radio-go-go/internal"
	mpv "github.com/kghose/radio-go-go/internal/mpv"
	radio_browser "github.com/kghose/radio-go-go/internal/radio_browser"
	//	"github.com/rivo/tview"
	"log"
	"log/slog"
	//	"time"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// We could have encapsulated all this as a struct, but instead
// we simply choose to use main.go as a struct like unit.
var server string

var searchResult []radio.Station
var history []radio.Station

var app *tview.Application

func setStationList(stations []radio.Station) {
	stationsListView.Clear()
	for _, station := range stations {
		stationsListView.AddItem(
			station.Details.Name, station.Details.URLResolved, 0, nil)
	}
}

var searchBarInputField *tview.InputField

func searchBarDone(key tcell.Key) {
	if key == tcell.KeyEnter {
		keywords := searchBarInputField.GetText()
		slog.Info(keywords)
		stations, err := radio_browser.StationSearch(keywords, server)
		if err != nil {
			slog.Error("Error searching for stations.")
		}
		searchResult = radio.SearchResults(stations, history)
		setStationList(searchResult)
		pages.SwitchToPage("Stations")
	}
	if key == tcell.KeyEsc {
		// Close the popup without doing anything
		pages.SwitchToPage("Stations")
	}
}

var pages *tview.Pages
var stationsListView *tview.List

func userKeyPress(event *tcell.EventKey) *tcell.EventKey {
	if event.Rune() == 'h' {
		setStationList(history)
		return nil
	}
	if event.Rune() == 's' {
		setStationList(searchResult)
		return nil
	}
	if event.Rune() == 'S' {
		pages.ShowPage("Search")
		return nil
	}
	if event.Rune() == 'q' {
		app.Stop()
		return nil
	}
	return event
}

func main() {
	mpv_player := mpv.Player{}
	mpv_player.Start()
	defer mpv_player.Quit()

	servers, err := radio_browser.GetAvailableServers()
	if err != nil {
		log.Fatal("Could not find radio browser servers")
	}

	server = radio_browser.PickRandomServer(servers)
	history, err = radio.LoadHistory()

	searchBarInputField = tview.NewInputField()
	searchBarInputField.SetFieldWidth(70).
		SetDoneFunc(searchBarDone)

	app = tview.NewApplication()

	pages = tview.NewPages()

	searchBar := tview.NewGrid().
		SetColumns(0, 80, 0).
		SetRows(0, 1, 0).
		AddItem(searchBarInputField, 1, 1, 1, 1, 0, 0, true)

	stationsListView = tview.NewList()

	stationsListView.SetSelectedFunc(
		func(_ int, _ string, url string, _ rune) {
			r := mpv_player.Play(url)
			slog.Info(r.Error)
		},
	)

	stationsListView.SetInputCapture(userKeyPress)

	pages.AddPage("Stations", stationsListView, true, true)
	pages.AddPage("Search", searchBar, true, true)

	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}

}

