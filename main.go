package main

import (
	//	"fmt"
	//	"github.com/gdamore/tcell/v2"
	radio "github.com/kghose/radio-go-go/internal"
	mpv "github.com/kghose/radio-go-go/internal/mpv"
	radio_browser "github.com/kghose/radio-go-go/internal/radio_browser"
	//	"github.com/rivo/tview"
	"log/slog"
	//	"time"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var searchBoxWidth = 50

type App struct {
	mpvPlayer mpv.Player

	server string

	searchResult []radio.Station
	history      []radio.Station

	ui                  *tview.Application
	stationsListView    *tview.List
	searchBarInputField *tview.InputField
	pages               *tview.Pages
}

func (app *App) setStationList(stations []radio.Station) {
	app.stationsListView.Clear()
	for _, station := range stations {
		app.stationsListView.AddItem(
			station.Details.Name, station.Details.URLResolved, 0, nil)
	}
}

func (app *App) searchBarDone(key tcell.Key) {
	if key == tcell.KeyEnter {
		keywords := app.searchBarInputField.GetText()
		stations, err := radio_browser.StationSearch(keywords, app.server)
		if err != nil {
			slog.Error("Error searching for stations.")
		}
		app.searchResult = radio.SearchResults(stations, app.history)
		app.setStationList(app.searchResult)
		app.pages.SwitchToPage("Stations")
	}
	if key == tcell.KeyEsc {
		// Close the popup without doing anything
		app.pages.SwitchToPage("Stations")
	}
}

func (app *App) playThis(_ int, _ string, url string, _ rune) {
	r := app.mpvPlayer.Play(url)
	slog.Info("Play", "url", url, "mpv", r.Error)
	app.history = radio.AddToHistory(url, app.searchResult, app.history)
}

func (app *App) favoriteThis(url string) {
	slog.Info("Fave", "url", url)
	app.history = radio.AddToFavorites(url, app.searchResult, app.history)
}

func (app *App) userKeyPress(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'h':
		app.setStationList(app.history)
	case 's':
		app.setStationList(app.searchResult)
	case 'S':
		app.pages.ShowPage("Search")
	case 'f':
		app.setStationList(radio.Favorites(app.history))
	case 'F':
		_, url := app.stationsListView.GetItemText(
			app.stationsListView.GetCurrentItem())
		app.favoriteThis(url)
	case 'p':
		app.mpvPlayer.TogglePause()
	case 'q':
		app.ui.Stop()
	default:
		return event
	}
	return nil
}


func main() {
	slogger, closeFunc := radio.SetupLoggingToFile()
	defer closeFunc()
	slog.SetDefault(slogger)

	app := App{}

	app.mpvPlayer = mpv.Player{}
	app.mpvPlayer.Start()
	defer app.mpvPlayer.Quit()

	servers, err := radio_browser.GetAvailableServers()
	if err != nil {
		slog.Error("Could not find radio browser servers")
	}

	app.server = radio_browser.PickRandomServer(servers)
	app.history, err = radio.LoadHistory()

	app.searchBarInputField = tview.NewInputField()
	app.searchBarInputField.SetFieldWidth(searchBoxWidth).
		SetDoneFunc(app.searchBarDone)

	app.ui = tview.NewApplication()

	app.pages = tview.NewPages()

	searchBar := tview.NewGrid().
		SetColumns(1, searchBoxWidth, 1).
		SetRows(1).
		SetBorders(true).
		AddItem(app.searchBarInputField, 0, 1, 1, 1, 0, 0, true)

	app.stationsListView = tview.NewList()
	app.stationsListView.SetSelectedFunc(app.playThis)
	app.stationsListView.SetInputCapture(app.userKeyPress)

	app.pages.AddPage("Stations", app.stationsListView, true, true)
	app.pages.AddPage("Search", searchBar, true, false)

	if err := app.ui.SetRoot(app.pages, true).Run(); err != nil {
		panic(err)
	}

}
