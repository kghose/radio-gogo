package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	radio "github.com/kghose/radio-go-go/internal"
	mpv "github.com/kghose/radio-go-go/internal/mpv"
	radio_browser "github.com/kghose/radio-go-go/internal/radio_browser"
)

var searchBoxWidth = 80

type ListMode int

const (
	HISTORY ListMode = iota
	SEARCH
	FAVES
)

type ListPos struct {
	topRow      int
	selectedRow int
}

type StationsListState struct {
	listPos  []ListPos
	listMode ListMode
}

func (ls *StationsListState) saveState(lv *tview.List) {
	offset, _ := lv.GetOffset()
	selRow := lv.GetCurrentItem()
	ls.listPos[ls.listMode] = ListPos{offset, selRow}
}

func (ls *StationsListState) loadState(newMode ListMode, lv *tview.List) {
	// First check if state is still suitable for the list
	if lv.GetItemCount() < ls.listPos[newMode].topRow {
		ls.listPos[newMode].topRow = lv.GetItemCount()
	}
	if lv.GetItemCount() < ls.listPos[newMode].selectedRow {
		ls.listPos[newMode].selectedRow = lv.GetItemCount()
	}

	lv.SetOffset(ls.listPos[newMode].topRow, 0)
	lv.SetCurrentItem(ls.listPos[newMode].selectedRow)

	ls.listMode = newMode
}

type App struct {
	mpvPlayer mpv.Player

	server string

	searchResult       []radio.Station
	history            []radio.Station
	lastSearchKeywords string

	ui                  *tview.Application
	stationsListView    *tview.List
	stationsListState   StationsListState
	searchBarInputField *tview.InputField
	nowPlayingBox       *tview.TextView
	pages               *tview.Pages
}

func stationViewTitleString(station *radio.Station) string {
	title := station.Details.Name
	if station.Favorite {
		title = "(*)" + title
	}
	return title
}

func (app *App) setStationList(stations []radio.Station) {
	app.stationsListView.Clear()
	for _, station := range stations {
		app.stationsListView.AddItem(
			stationViewTitleString(&station), station.Details.URLResolved, 0, nil)
	}
	app.stationsListView.SetCurrentItem(0)
}

func (app *App) resetSearchCursor() {
	app.stationsListState.listPos[SEARCH] = ListPos{}
	if app.stationsListState.listMode != SEARCH {
		return
	}
	app.stationsListView.SetOffset(0, 0)
	app.stationsListView.SetCurrentItem(0)
}

func (app *App) setListTo(t ListMode) {

	list := []radio.Station{}
	title := ""

	switch t {
	case HISTORY:
		title = "History"
		list = app.history
	case SEARCH:
		if app.lastSearchKeywords != "" {
			title = app.lastSearchKeywords
		} else {
			title = "Search"
		}
		list = app.searchResult
	case FAVES:
		title = "Faves"
		list = radio.Favorites(app.history)
	}

	app.stationsListState.saveState(app.stationsListView)
	app.setStationList(list)
	app.stationsListState.loadState(t, app.stationsListView)
	app.stationsListView.SetTitle(fmt.Sprintf("%s (%d)", title, len(list)))
	app.pages.SwitchToPage("Stations")
	app.ui.SetFocus(app.stationsListView)
}

func (app *App) searchBarDone(key tcell.Key) {
	if key == tcell.KeyEnter {
		keywords := app.searchBarInputField.GetText()
		stations, err := radio_browser.StationSearch(keywords, app.server)
		if err != nil {
			slog.Error("Error searching for stations.")
		}
		app.searchResult = radio.SearchResults(stations, app.history)
		app.lastSearchKeywords = keywords
		app.resetSearchCursor()
		app.setListTo(SEARCH)
	}
	if key == tcell.KeyEsc {
		// Close the popup without doing anything
		app.pages.SwitchToPage("Stations")
		app.ui.SetFocus(app.stationsListView)
	}
}

func (app *App) playThis(_ int, _ string, url string, _ rune) {
	r := app.mpvPlayer.Play(url)
	slog.Info("Play", "url", url, "mpv", r.Error)
	app.history = radio.AddToHistory(url, app.searchResult, app.history)
	SaveHistory(app.history)
}

func (app *App) updateNowPlayingBox() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	shs := NewSongHistorySaver()

	for range ticker.C {
		meta := app.mpvPlayer.Meta()
		text := fmt.Sprintf(
			"Station: %s\nSummary: %s\nGenre: %s\nTrack: %s",
			meta.Name, meta.Description, meta.Genre, meta.Title)
		app.ui.QueueUpdateDraw(func() { app.nowPlayingBox.SetText(text) })
		shs.save(meta.Title)
	}
}

func (app *App) favoriteThis(url string) {
	slog.Info("Fave", "url", url)
	app.history = radio.AddToFavorites(url, app.searchResult, app.history)
	SaveHistory(app.history)
}

func (app *App) userKeyPress(event *tcell.EventKey) *tcell.EventKey {
	if app.searchBarInputField.HasFocus() {
		return event
	}

	switch event.Rune() {
	case 'h':
		app.setListTo(HISTORY)
	case 's':
		app.setListTo(SEARCH)
	case '/':
		app.pages.ShowPage("Search")
		app.ui.SetFocus(app.searchBarInputField)
	case 'f':
		app.setListTo(FAVES)
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
	slogger, closeFunc := SetupLoggingToFile()
	defer closeFunc()
	slog.SetDefault(slogger)

	app := App{}

	app.mpvPlayer = mpv.Player{}
	app.mpvPlayer.Start()
	defer app.mpvPlayer.Quit()

	servers, err := radio_browser.GetAvailableServers()
	if err != nil {
		slog.Error("Could not find radio browser servers")
		servers = []string{""}
	}

	app.server = radio_browser.PickRandomServer(servers)
	app.history, err = LoadHistory()

	app.searchBarInputField = tview.NewInputField()
	app.searchBarInputField.SetFieldWidth(searchBoxWidth).
		SetDoneFunc(app.searchBarDone)

	app.ui = tview.NewApplication()

	app.pages = tview.NewPages()

	searchBar := tview.NewGrid().
		SetColumns(0, searchBoxWidth, 0).
		SetRows(6, 1).
		SetBorders(true).
		AddItem(app.searchBarInputField, 1, 1, 1, 1, 0, 0, true)

	app.nowPlayingBox = tview.NewTextView()
	app.nowPlayingBox.
		SetBorder(true).
		SetTitleAlign(tview.AlignRight).
		SetTitle("Playing")
	go app.updateNowPlayingBox()

	app.stationsListView = tview.NewList()
	app.stationsListView.SetSelectedFunc(app.playThis)
	app.stationsListView.SetTitleAlign(tview.AlignRight).SetBorder(true)
	app.stationsListState.listPos = []ListPos{{}, {}, {}}

	mainGrid := tview.NewGrid().
		SetColumns(100).
		SetRows(6, 0).
		AddItem(app.nowPlayingBox, 0, 0, 1, 1, 5, 80, false).
		AddItem(app.stationsListView, 1, 0, 1, 1, 20, 80, true)

	app.pages.AddPage("Stations", mainGrid, true, true)
	app.pages.AddPage("Search", searchBar, true, false)

	app.ui.SetInputCapture(app.userKeyPress)

	app.setListTo(HISTORY)

	if err := app.ui.SetRoot(app.pages, true).Run(); err != nil {
		panic(err)
	}
}
