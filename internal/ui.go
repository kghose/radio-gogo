/*
Handle some of the UI elements.
*/
package radio

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/kghose/radio-go-go/internal/mpv"
)

const searchBarWidth = 80

type PageName string

const (
	histPage   PageName = "History"
	searchPage PageName = "Search Results"
	favesPage  PageName = "Faves"
)

type StationsView struct {
	pages *tview.Pages
	lists map[PageName]*tview.List
	title map[PageName]string
}

func (sv *StationsView) currentPage() PageName {
	page, _ := sv.pages.GetFrontPage()
	return PageName(page)
}

func (sv *StationsView) setup(playThis func(int, string, string, rune)) {
	sv.pages = tview.NewPages()
	sv.lists = make(map[PageName]*tview.List)
	sv.title = make(map[PageName]string)

	for _, page := range []PageName{histPage, searchPage, favesPage} {
		l := tview.NewList()
		l.
			ShowSecondaryText(false).
			SetSelectedFunc(playThis).
			SetTitleAlign(tview.AlignRight).
			SetBorder(false)
		sv.lists[page] = l
		sv.pages.AddPage(string(page), l, true, false)
	}

}

func itemTitle(station *Station) string {
	heart := ""
	if station.Favorite {
		heart = "[red]\u2764[-]"
	}
	return fmt.Sprintf(
		"%-1s %-30.30s [blue]%s",
		heart, station.Name, station.URLResolved)
}

func (sv *StationsView) set(stations []*Station, pageName PageName, title string, reset_view bool) {
	selRow := 0
	offRow := 0

	list := sv.lists[pageName]

	if !reset_view {
		selRow = list.GetCurrentItem()
		offRow, _ = list.GetOffset()
	}

	list.Clear()
	for _, station := range stations {
		list.AddItem(
			itemTitle(station), station.URLResolved, 0, nil)
	}
	list.SetCurrentItem(selRow)
	list.SetOffset(offRow, 0)
	sv.title[pageName] = fmt.Sprintf("%s (%d)", title, len(stations))
}

type ViewName string

const (
	mainView        ViewName = "Main"
	searchBarPopup  ViewName = "Search Popup"
	helpPopup       ViewName = "Help Popup"
	stationsPage    ViewName = "Stations"
	playedsongsPage ViewName = "Played songs"
)

type UI struct {
	app         *tview.Application
	pages       *tview.Pages
	stnInfoPane *tview.TextView
	listsPane   *tview.Pages
	searchBar   *tview.InputField
	playedsongs *tview.TextView

	stationsView StationsView
}

func (ui *UI) ShowSearchBar() {
	ui.pages.ShowPage(string(searchBarPopup))
	ui.searchBar.SetText("")
	ui.app.SetFocus(ui.searchBar)
}

func (ui *UI) HideSearchBar() {
	ui.pages.HidePage(string(searchBarPopup))
	ui.app.SetFocus(ui.stationsView.pages)
}

func (ui *UI) ShowHelp() {
	ui.pages.ShowPage(string(helpPopup))
}

func (ui *UI) HideHelp() {
	ui.pages.HidePage(string(helpPopup))
}

func (ui *UI) RefreshPlayedsongs(songs iter.Seq[string]) {
	ui.app.QueueUpdateDraw(func() {
		ui.playedsongs.SetText(strings.Join(slices.Collect(songs), "\n"))
	})
}

var playStateString = map[bool]string{
	true:  "[green]\u25b6[-]",
	false: "[red]\u23f8[-]",
}

func (ui *UI) SetNowPlaying(meta mpv.MpvMetadata) {
	text := fmt.Sprintf(
		`[red]Station: %s
[white]Summary: %s
[blue]Genre: %s
[yellow]Track: %s %s`,
		meta.Name,
		meta.Description,
		meta.Genre,
		playStateString[meta.Playing], meta.Title)
	ui.app.QueueUpdateDraw(func() { ui.stnInfoPane.SetText(text) })
}

func (ui *UI) show(pageName PageName) {
	ui.listsPane.SwitchToPage(string(stationsPage))
	ui.stationsView.pages.SwitchToPage(string(pageName))
}

func (ui *UI) ShowHist() {
	ui.show(histPage)
}

func (ui *UI) ShowFaves() {
	ui.show(favesPage)
}

func (ui *UI) ShowSearch() {
	ui.show(searchPage)
}

func (ui *UI) ShowPlayedSongs() {
	ui.listsPane.SwitchToPage(string(playedsongsPage))
}

func (ui *UI) SelectedURL() string {
	cp := ui.stationsView.currentPage()
	idx := ui.stationsView.lists[cp].GetCurrentItem()
	_, url := ui.stationsView.lists[cp].GetItemText(idx)
	return url
}

func (ui *UI) SetHist(stations []*Station) {
	ui.stationsView.set(
		stations,
		histPage, string(histPage), false)
}

func (ui *UI) SetFaves(stations []*Station) {
	ui.stationsView.set(
		stations,
		favesPage, string(favesPage), false)
}

func (ui *UI) SetSearch(stations []*Station, keywords string) {
	ui.stationsView.set(
		stations,
		searchPage, keywords, false)
}

func (ui *UI) ResetSearchScroll() {
	ui.stationsView.lists[searchPage].SetCurrentItem(0)
	ui.stationsView.lists[searchPage].SetOffset(0, 0)
}

func (ui *UI) RefreshLists(index map[string]*Station, keywords string) {
	ui.SetHist(SortLastPlayed(History(index)))
	ui.SetFaves(SortAlpha(Faves(index)))
	ui.SetSearch(SortAlpha(Search(index)), keywords)
}

type KeyFunc struct {
	Help string
	Fn   func()
}

func makeHelpText(keyMap map[rune]KeyFunc) string {
	text := "\n\n[yellow]Keys[-]\n\n"
	for _, k := range slices.Sorted(maps.Keys(keyMap)) {
		text = text +
			fmt.Sprintf(
				"[white]%s [red]: [green]%s\n",
				string(k), keyMap[k].Help)
	}
	text = text + searchHelp
	return text
}

func (ui *UI) Setup(
	keyMap map[rune]KeyFunc,
	searchFunc func(string),
	playFunc func(int, string, string, rune)) {

	ui.stnInfoPane = tview.NewTextView()
	ui.stnInfoPane.SetBorder(false)
	ui.stnInfoPane.SetDynamicColors(true)

	ui.stationsView.setup(playFunc)

	ui.playedsongs = tview.NewTextView().
		SetSize(0, 80).
		SetDynamicColors(true)

	ui.listsPane = tview.NewPages()
	ui.listsPane.AddPage(string(stationsPage), ui.stationsView.pages, true, true)
	ui.listsPane.AddPage(string(playedsongsPage), ui.playedsongs, true, false)
	ui.listsPane.SetDrawFunc(
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
			title := ui.stationsView.title[ui.stationsView.currentPage()]
			name, _ := ui.listsPane.GetFrontPage()
			if name == string(playedsongsPage) {
				title = name
			}
			tview.Print(
				screen, " "+title,
				x, y, width, tview.AlignRight, tcell.ColorYellow)

			// Return the inner rectangle where content should be drawn
			// (We subtract 1 from the top to account for the title line)
			return x, y + 1, width, height - 1
		})

	mainPageGrid := tview.NewGrid().
		SetColumns(100).
		SetRows(4, 0).
		AddItem(ui.stnInfoPane, 0, 0, 1, 1, 3, 80, false).
		AddItem(ui.listsPane, 1, 0, 1, 1, 20, 80, true)

	ui.searchBar = tview.NewInputField().
		SetFieldWidth(searchBarWidth).
		SetDoneFunc(func(k tcell.Key) {
			if k == tcell.KeyEnter {
				searchFunc(ui.searchBar.GetText())
			}
			ui.HideSearchBar()
		}).
		SetFieldBackgroundColor(tcell.GetColor("white")).
		SetFieldTextColor(tcell.GetColor("black"))

	searchBarGrid := tview.NewGrid().
		SetColumns(0, searchBarWidth, 0).
		SetRows(2, 1).
		SetBorders(true).
		AddItem(ui.searchBar, 1, 1, 1, 1, 0, 0, true)

	helpText := tview.NewTextView().
		SetSize(0, 80).
		SetDynamicColors(true).
		SetRegions(true).
		SetDoneFunc(func(_ tcell.Key) {
			ui.HideHelp()
		}).
		SetLabel("Help").
		SetText(makeHelpText(keyMap))
	helpText.SetBorder(true)

	ui.pages = tview.NewPages()
	ui.pages.AddPage(string(mainView), mainPageGrid, true, true)
	ui.pages.AddPage(string(searchBarPopup), searchBarGrid, true, false)
	ui.pages.AddPage(string(helpPopup), helpText, true, false)

	ui.app = tview.NewApplication()
	ui.app.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if ui.searchBar.HasFocus() {
			return e
		}
		if keyPress, ok := keyMap[e.Rune()]; ok {
			keyPress.Fn()
			return nil
		} else {
			return e
		}
	})
	ui.app.SetRoot(ui.pages, true)
}

func (ui *UI) Run() error {
	return ui.app.Run()
}

func (ui *UI) Stop() {
	ui.app.Stop()
}
