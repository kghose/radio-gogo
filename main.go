package main

import (
	"log/slog"
	"time"

	radio "github.com/kghose/radio-go-go/internal"
	mpv "github.com/kghose/radio-go-go/internal/mpv"
	radio_browser "github.com/kghose/radio-go-go/internal/radio_browser"
)

/*
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

	func itemTitle(station *radio.Station) string {
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
				itemTitle(&station), station.Details.URLResolved, 0, nil)
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

	func (app *App) hideSearchBar() {
		app.stationsPane.widget.HidePage(SEARCH_BAR)
	}

	func (app *App) searchBarText() string {
		return app.stationsPane.searchBarInputField.GetText()
	}
*/
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
	ui.Setup(
		func(r rune) bool {
			switch r {
			case 'h':
				ui.ShowHist()
			case 's':
				ui.ShowSearch()
			case '/':
				ui.ShowSearchBar()
			case 'f':
				ui.ShowFaves()
			case '=':
				radio.UpdateIndex(ui.SelectedURL(), stationIndex, radio.FAVE)
				ui.RefreshLists(stationIndex, searchUrls, keywords)
				SaveHistory(stationIndex)
			case '-':
				radio.UpdateIndex(ui.SelectedURL(), stationIndex, radio.UNFAVE)
				ui.RefreshLists(stationIndex, searchUrls, keywords)
				SaveHistory(stationIndex)
			case 'p':
				mpvPlayer.TogglePause()
			case 'q':
				ui.Stop()
			default:
				return false
			}
			return true
		},
		func(kw string) {
			keywords = kw
			stations, err := radio_browser.StationSearch(keywords, server)
			if err != nil {
				slog.Error("Error searching for stations.")
			}
			stationIndex, searchUrls =
				radio.MakeNewIndexFromSearch(stations, stationIndex)
			ui.RefreshLists(stationIndex, searchUrls, keywords)
			ui.ShowSearch()
		},
		func(idx int, _ string, url string, _ rune) {
			r := mpvPlayer.Play(url)
			slog.Info("Play", "url", url, "mpv", r.Error)
			radio.UpdateIndex(url, stationIndex, radio.PLAYED)
			ui.RefreshLists(stationIndex, searchUrls, keywords)
			SaveHistory(stationIndex)
		},
	)

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for ; ; <-ticker.C {
			meta := mpvPlayer.Meta()
			ui.SetNowPlaying(meta)
			shs.save(meta.Title)
		}
	}()

	ui.RefreshLists(stationIndex, searchUrls, keywords)
	ui.ShowHist()
	if err := ui.Run(); err != nil {
		panic(err)
	}
}
