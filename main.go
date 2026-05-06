package main

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	radio "github.com/kghose/radio-go-go/internal"
	mpv "github.com/kghose/radio-go-go/internal/mpv"
	radio_browser "github.com/kghose/radio-go-go/internal/radio_browser"
)

var searchBoxWidth = 80

const STATION_NAME_JUNK_CHARS = ".-+*# "

type ListName string

const (
	HISTORY        ListName = "History"
	SEARCH_RESULTS ListName = "Search Results"
	FAVES          ListName = "Faves"
)

type PopupName string

const SEARCH_BAR string = "Search Bar"

type StationsList struct {
	widget   *tview.List
	title    string
	stations []radio.Station
}

func stationString(station *radio.Station) string {
	name := strings.TrimLeft(station.Details.Name, STATION_NAME_JUNK_CHARS)
	if station.Favorite {
		name = "❤️" + name
	}
	return fmt.Sprintf("%-*.*s [blue] %s", 30, 30, name, station.Details.URLResolved)
}

func (sv *StationsList) set(stations []radio.Station, title string, reset_view bool) {
	selRow := 0
	offRow := 0

	if !reset_view {
		selRow = sv.widget.GetCurrentItem()
		offRow, _ = sv.widget.GetOffset()
	}

	sv.widget.Clear()
	for _, station := range stations {
		sv.widget.AddItem(
			stationString(&station), station.Details.URLResolved, 0, nil)
	}
	sv.widget.SetCurrentItem(selRow)
	sv.widget.SetOffset(offRow, 0)
	sv.title = fmt.Sprintf("%s (%d)", title, len(stations))
	sv.stations = stations
}

type StationsPane struct {
	widget              *tview.Pages
	lists               map[ListName]*StationsList
	currentList         ListName
	searchBarInputField *tview.InputField
	searchBar           *tview.Grid
}

func (pane *StationsPane) setupWidget() {
	pane.widget = tview.NewPages()
	pane.widget.SetDrawFunc(
		// Custom border
		func(
			screen tcell.Screen,
			x, y, width, height int) (int, int, int, int) {
			// Line
			for cx := x; cx < x+width; cx++ {
				tview.Print(
					screen,
					string(tview.BoxDrawingsLightHorizontal),
					cx, y, 1, tview.AlignCenter,
					tcell.ColorWhite)
			}
			// Title
			tview.Print(
				screen, " "+pane.lists[pane.currentList].title,
				x, y, width, tview.AlignRight, tcell.ColorYellow)

			// Return the inner rectangle where content should be drawn
			// (We subtract 1 from the top to account for the title line)
			return x, y + 1, width, height - 1
		})
}

func (sp *StationsPane) setupSearchBar(searchFunc func(tcell.Key)) {
	sp.searchBarInputField = tview.NewInputField()
	sp.searchBarInputField.
		SetFieldWidth(searchBoxWidth).
		SetDoneFunc(searchFunc).
		SetFieldBackgroundColor(tcell.GetColor("white")).
		SetFieldTextColor(tcell.GetColor("black"))
	sp.searchBar = tview.NewGrid().
		SetColumns(0, searchBoxWidth, 0).
		SetRows(2, 1).
		SetBorders(true).
		AddItem(sp.searchBarInputField, 1, 1, 1, 1, 0, 0, true)
	sp.widget.AddPage(SEARCH_BAR, sp.searchBar, true, false)
}

func (sp *StationsPane) setup(
	playThis func(int, string, string, rune),
	searchFunc func(tcell.Key)) {
	sp.setupWidget()
	sp.lists = make(map[ListName]*StationsList)
	for _, page := range []ListName{HISTORY, SEARCH_RESULTS, FAVES} {
		tv := tview.NewList()
		sp.lists[page] = &StationsList{tv, string(page), []radio.Station{}}
		tv.
			ShowSecondaryText(false).
			SetSelectedFunc(playThis).
			SetTitleAlign(tview.AlignRight).
			SetBorder(false)

		sp.widget.AddPage(string(page), tv, true, false)
	}
	sp.setupSearchBar(searchFunc)
}

func (sp *StationsPane) switchTo(listName ListName) {
	sp.widget.SwitchToPage(string(listName))
	sp.currentList = listName
}

type App struct {
	mpvPlayer mpv.Player

	server string

	searchResult       []radio.Station
	history            []radio.Station
	lastSearchKeywords string

	ui            *tview.Application
	stationsPane  StationsPane
	nowPlayingBox *tview.TextView
}

func (app *App) loadHist() {
	stations, err := LoadHistory()
	// TODO: Handle errors
	if err != nil {
	}
	app.stationsPane.lists[HISTORY].stations = stations
	app.updateLists()
}

func (app *App) getHist() *[]radio.Station {
	return &app.stationsPane.lists[HISTORY].stations
}

func (app *App) saveHist() {
	SaveHistory(*app.getHist())
}

func (app *App) currentStations() *[]radio.Station {
	return &app.stationsPane.lists[app.stationsPane.currentList].stations
}

func (app *App) updateLists() {
	app.
		stationsPane.
		lists[HISTORY].set(
		*app.getHist(),
		string(HISTORY),
		false,
	)
	app.
		stationsPane.
		lists[FAVES].set(
		radio.Favorites(*app.getHist()),
		string(FAVES),
		false,
	)
}

func (app *App) updateHist(idx int, op radio.StationOp) {
	station := (*app.currentStations())[idx]
	slog.Info(string(op), "url", station.Details.URLResolved)
	radio.UpdateHist(
		station,
		app.getHist(),
		op)
	app.saveHist()
	app.updateLists()
}

func (app *App) playThis(idx int, _ string, url string, _ rune) {
	r := app.mpvPlayer.Play(url)
	slog.Info("Play", "url", url, "mpv", r.Error)
	app.updateHist(idx, radio.PLAYED)
}

func (app *App) currentStationIdx() int {
	return app.
		stationsPane.
		lists[app.stationsPane.currentList].
		widget.
		GetCurrentItem()
}

func (app *App) faveThis() {
	app.updateHist(app.currentStationIdx(), radio.FAVE)
	app.saveHist()
	app.updateLists()
}

func (app *App) unfaveThis() {
	app.updateHist(app.currentStationIdx(), radio.UNFAVE)
	app.saveHist()
	app.updateLists()
}


func (app *App) showSearchBar() {
	app.stationsPane.widget.ShowPage(SEARCH_BAR)
	app.stationsPane.searchBarInputField.SetText("")
	app.ui.SetFocus(app.stationsPane.searchBarInputField)
}

func (app *App) hideSearchBar() {
	app.stationsPane.widget.HidePage(SEARCH_BAR)
}

func (app *App) searchBarText() string {
	return app.stationsPane.searchBarInputField.GetText()
}

func (app *App) searchFor(key tcell.Key) {
	if key == tcell.KeyEnter {
		keywords := app.searchBarText()
		stations, err := radio_browser.StationSearch(keywords, app.server)
		if err != nil {
			slog.Error("Error searching for stations.")
		}
		app.
			stationsPane.
			lists[SEARCH_RESULTS].set(
			radio.SearchResults(stations, *app.getHist()),
			keywords,
			true,
		)
		app.stationsPane.switchTo(SEARCH_RESULTS)
	}
	if key == tcell.KeyEsc {
		// Close the popup without doing anything
		app.hideSearchBar()
	}
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

func (app *App) userKeyPress(event *tcell.EventKey) *tcell.EventKey {
	if app.stationsPane.searchBarInputField.HasFocus() {
		return event
	}

	switch event.Rune() {
	case 'h':
		app.stationsPane.switchTo(HISTORY)
	case 's':
		app.stationsPane.switchTo(SEARCH_RESULTS)
	case '/':
		app.showSearchBar()
	case 'f':
		app.stationsPane.switchTo(FAVES)
	case '=':
		app.faveThis()
	case '-':
		app.unfaveThis()
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

	app.stationsPane.setup(app.playThis, app.searchFor)
	app.loadHist()

	app.ui = tview.NewApplication()

	app.nowPlayingBox = tview.NewTextView()
	app.nowPlayingBox.
		SetBorder(false).
		SetTitleAlign(tview.AlignRight)
	go app.updateNowPlayingBox()

	mainGrid := tview.NewGrid().
		SetColumns(100).
		SetRows(4, 0).
		AddItem(app.nowPlayingBox, 0, 0, 1, 1, 3, 80, false).
		AddItem(app.stationsPane.widget, 1, 0, 1, 1, 20, 80, true)

	app.ui.SetInputCapture(app.userKeyPress)

	app.stationsPane.switchTo(HISTORY)

	if err := app.ui.SetRoot(mainGrid, true).Run(); err != nil {
		panic(err)
	}
}
