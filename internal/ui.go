package radio

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const STATION_METADATA_REFRESH_INTERVAL = time.Second

const (
	SEARCH_PANE = 0
	HIST_PANE   = 1
	FAV_PANE    = 2
)

type PaneState struct {
	pane       int
	pane_index []int
}

type RadioUI struct {
	device       Radio
	player       Player
	pane_state   PaneState
	app          *tview.Application
	search_bar   *tview.InputField
	tab_title    *tview.TextView
	station_list *tview.List
	now_playing  *tview.TextView
	status_bar   *tview.TextView
}

func (r *RadioUI) Run() {

	r.pane_state.pane_index = []int{0, 0, 0}
	r.device = NewRadio()
	if err := r.device.Load_user_data(); err != nil {
		panic(err)
	}

	r.player.Start()
	r.app = tview.NewApplication()
	r.setup_UI(r.app)
	go r.RefreshServers()
	go r.periodically_update_stream_metadata()
	if err := r.app.Run(); err != nil {
		r.app.Stop()
	}
	r.player.Quit()

	if err := r.device.Save_user_data(); err != nil {
		panic(err)
	}

}

func (r *RadioUI) setup_UI(app *tview.Application) {

	r.search_bar = tview.NewInputField().
		SetDoneFunc(r.Search)
	r.tab_title = tview.NewTextView()
	r.tab_title.SetTextAlign(tview.AlignCenter).
		SetBackgroundColor(tcell.ColorDarkBlue)
	r.tab_title.SetTextColor(tcell.ColorWhiteSmoke)
	r.station_list = tview.NewList().
		ShowSecondaryText(true).
		SetSelectedFunc(func(_ int, _ string, url string, _ rune) {
			r.play(url)
		})
	r.now_playing = tview.NewTextView()
	r.status_bar = tview.NewTextView()

	grid := tview.NewGrid().
		SetRows(1, -1, 4, 1).
		SetColumns(10, -1).
		SetBorders(true).
		SetBordersColor(tcell.ColorGreenYellow)
	grid.AddItem(r.search_bar, 0, 1, 1, 1, 0, 0, true)
	grid.AddItem(r.tab_title, 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(r.station_list, 1, 0, 1, 2, 0, 0, false)
	grid.AddItem(r.now_playing, 2, 0, 1, 2, 0, 0, false)
	grid.AddItem(r.status_bar, 3, 0, 1, 2, 0, 0, false)
	app.SetRoot(grid, true).SetFocus(grid).SetInputCapture(r.input_capture)
}

func (r *RadioUI) RefreshServers() {
	r.app.QueueUpdateDraw(func() {
		r.status_bar.SetText("Refreshing server list ...")
	})
	err := r.device.Refresh_servers()
	for err != nil {
		r.app.QueueUpdateDraw(
			func() {
				r.status_bar.SetText(
					fmt.Sprintf("Error refreshing server list %s",
						err.Error()))
			})
		time.Sleep(1 * time.Second)
		err = r.device.Refresh_servers()
	}
	r.app.QueueUpdateDraw(
		func() { r.status_bar.SetText(fmt.Sprintf("Found %d Radio Browser servers.", len(r.device.Servers))) },
	)
}

func (r *RadioUI) Search(key tcell.Key) {
	if key != tcell.KeyEnter {
		return
	}

	// Wait until servers have been found ...
	if len(r.device.Servers) == 0 {
		return
	}

	go func() {
		var err error
		err = r.device.FindByTag([]string{r.search_bar.GetText()})
		r.app.QueueUpdateDraw(func() {

			if err != nil {
				r.status_bar.SetText(
					fmt.Sprintf("Search error: %s", err.Error()))
				return
			}

			r.status_bar.SetText(
				fmt.Sprintf("Found %d stations.",
					r.device.Stations.Len()))

			r.show_search()
		})
	}()

	r.status_bar.SetText("Searching ...")

}

func (r *RadioUI) update_station_list(stations *StationSet) {
	r.station_list.Clear()
	for i := range stations.Stations {
		r.station_list.AddItem(stations.Stations[i].Name, stations.Stations[i].Url, 0, nil)
	}
	r.station_list.SetCurrentItem(0)
	if r.station_list.GetItemCount() > 0 {
		r.app.SetFocus(r.station_list)
	}
}

func (r *RadioUI) show_search() {
	r.pane_state.pane_index[r.pane_state.pane] = r.station_list.GetCurrentItem()
	r.tab_title.SetText("Search")
	r.pane_state.pane = SEARCH_PANE
	r.update_station_list(
		r.device.Stations)
	r.station_list.SetCurrentItem(r.pane_state.pane_index[r.pane_state.pane])
}

func (r *RadioUI) show_history() {
	r.pane_state.pane_index[r.pane_state.pane] = r.station_list.GetCurrentItem()
	r.tab_title.SetText("History")
	r.pane_state.pane = HIST_PANE
	r.update_station_list(
		r.device.User_data.Station_history)
	r.station_list.SetCurrentItem(r.pane_state.pane_index[r.pane_state.pane])
}

func (r *RadioUI) show_favorites() {
	r.pane_state.pane_index[r.pane_state.pane] = r.station_list.GetCurrentItem()
	r.tab_title.SetText("Favorites")
	r.pane_state.pane = FAV_PANE
	r.update_station_list(
		r.device.User_data.Station_favorites)
	r.station_list.SetCurrentItem(r.pane_state.pane_index[r.pane_state.pane])
}

func (r *RadioUI) play(url string) {
	station := r.device.Stations.By_url(url)
	resp := r.player.Play(station.Url)
	r.status_bar.SetText(resp.Error)
	r.device.Now_playing(station)
	r.now_playing.SetText(station.Name)
}

func (r *RadioUI) periodically_update_stream_metadata() {
	for {
		time.Sleep(STATION_METADATA_REFRESH_INTERVAL)
		if !r.player.playing {
			continue
		}
		r.app.QueueUpdateDraw(r.update_stream_metadata)
	}
}

func (r *RadioUI) update_stream_metadata() {
	meta := r.player.Meta()
	r.now_playing.SetText(
		fmt.Sprintf(
			"Station: %s\nSummary: %s\nGenre: %s\nTrack: %s",
			meta.Name,
			meta.Description,
			meta.Genre,
			meta.Title,
		),
	)
}

func (r *RadioUI) input_capture(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyTab {
		if r.app.GetFocus() == r.search_bar {
			r.app.SetFocus(r.station_list)
		} else {
			r.app.SetFocus(r.search_bar)
		}
		return nil

	}

	if r.app.GetFocus() != r.search_bar {
		switch event.Rune() {
		case 'q':
			r.app.Stop()
		case 'p':
			r.player.Pause()
		case 's':
			r.show_search()
		case 'h':
			r.show_history()
		case 'f':
			r.show_favorites()
			//case '=':
			//	r.add_to_favorites()
		}
	}
	return event
}
