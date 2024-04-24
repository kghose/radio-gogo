package radio

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const STATION_METADATA_REFRESH_INTERVAL = time.Second

type UIState struct {
	current_pane int
	pane_state   [3]PaneState
}

type PaneState struct {
	name   string
	offset int
	index  int
}

var pane_name = [3]string{"Search", "History", "Favorites"}

type RadioUI struct {
	device       Radio
	player       Player
	ui_state     UIState
	app          *tview.Application
	grid         *tview.Grid
	search_bar   *tview.InputField
	tab_title    *tview.TextView
	station_list *tview.List
	now_playing  *tview.TextView
	status_bar   *tview.TextView
	help_modal   *tview.Modal
	keymap       KeyMap
}

func (r *RadioUI) Run() {

	r.device = NewRadio()
	if err := r.device.Load_user_data(); err != nil {
		panic(err)
	}

	r.player.Start()
	r.app = tview.NewApplication()
	r.setup_UI()
	r.setup_keymap()

	go r.RefreshServers()
	go r.periodically_update_stream_metadata()

	go r.app.QueueUpdateDraw(func() {
		r.update_station_list(r.ui_state.current_pane)
		r.app.SetFocus(r.search_bar)
	})
	if err := r.app.Run(); err != nil {
		r.app.Stop()
	}
	r.player.Quit()

	if err := r.device.Save_user_data(); err != nil {
		panic(err)
	}

}

func (r *RadioUI) setup_UI() {

	r.ui_state.current_pane = STATION_LIST_HIST
	r.ui_state.pane_state[STATION_LIST_SEARCH].name = "Search"
	r.ui_state.pane_state[STATION_LIST_HIST].name = "History"
	r.ui_state.pane_state[STATION_LIST_FAV].name = "Favorites"

	r.search_bar = tview.NewInputField().
		SetDoneFunc(r.Search)
	r.tab_title = tview.NewTextView()
	r.tab_title.SetTextAlign(tview.AlignCenter).
		SetBackgroundColor(tcell.ColorDarkBlue)
	r.tab_title.SetTextColor(tcell.ColorWhiteSmoke)
	r.station_list = tview.NewList().
		ShowSecondaryText(false).
		SetSelectedFunc(func(idx int, _ string, _ string, _ rune) {
			r.play(idx)
		})
	r.now_playing = tview.NewTextView()
	r.status_bar = tview.NewTextView()
	r.help_modal = tview.NewModal()

	r.grid = tview.NewGrid().
		SetRows(1, -1, 4, 1).
		SetColumns(10, -1).
		SetBorders(true).
		SetBordersColor(tcell.ColorGreenYellow)
	r.grid.AddItem(r.search_bar, 0, 1, 1, 1, 0, 0, true)
	r.grid.AddItem(r.tab_title, 0, 0, 1, 1, 0, 0, false)
	r.grid.AddItem(r.station_list, 1, 0, 1, 2, 0, 0, false)
	r.grid.AddItem(r.now_playing, 2, 0, 1, 2, 0, 0, false)
	r.grid.AddItem(r.status_bar, 3, 0, 1, 2, 0, 0, false)
	r.app.SetRoot(r.grid, true).SetFocus(r.grid).SetInputCapture(r.input_capture)

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
					r.device.StationLists[STATION_LIST_SEARCH].Len()))

			r.update_station_list(STATION_LIST_SEARCH)
		})
	}()

	r.status_bar.SetText("Searching ...")

}

func (r *RadioUI) update_station_list(list int) {
	r.ui_state.pane_state[r.ui_state.current_pane].index = r.station_list.GetCurrentItem()
	r.ui_state.pane_state[r.ui_state.current_pane].offset, _ = r.station_list.GetOffset()

	r.ui_state.current_pane = list
	r.tab_title.SetText(r.ui_state.pane_state[list].name)
	stations := r.device.StationLists[list]

	r.station_list.Clear()
	for i := range stations.Stations {
		r.station_list.AddItem(
			fmt.Sprintf(
				"%s [gray](%s)",
				stations.Stations[i].Name,
				stations.Stations[i].Url),
			"",
			0,
			nil,
		)
	}
	if r.station_list.GetItemCount() > 0 {
		r.app.SetFocus(r.station_list)
		r.station_list.SetCurrentItem(r.ui_state.pane_state[list].index)
		r.station_list.SetOffset(r.ui_state.pane_state[list].offset, 0)
	}
}

func (r *RadioUI) favorite() {
	idx := r.station_list.GetCurrentItem()
	r.device.StationLists[STATION_LIST_FAV].add(
		r.device.StationLists[r.ui_state.current_pane].Stations[idx],
	)
}

func (r *RadioUI) remove() {
	idx := r.station_list.GetCurrentItem()
	r.device.StationLists[r.ui_state.current_pane].remove(
		r.device.StationLists[r.ui_state.current_pane].Stations[idx],
	)
	r.update_station_list(r.ui_state.current_pane)
}

func (r *RadioUI) play(idx int) {
	station := r.device.StationLists[r.ui_state.current_pane].Stations[idx]
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

	if r.app.GetFocus() == r.station_list {
		if mapping, ok := r.keymap[event.Rune()]; ok {
			mapping.fn()
		} else {
			//r.show_help()
		}
	}
	return event
}

type KeyMap map[rune]KeyMapping

type KeyMapping struct {
	fn   func()
	help string
}

func (r *RadioUI) setup_keymap() {
	r.keymap = KeyMap{
		'q': KeyMapping{
			fn:   r.app.Stop,
			help: "Quit",
		},
		'p': KeyMapping{
			fn:   func() { r.player.Pause() },
			help: "Pause/Unpause"},
		's': KeyMapping{
			fn:   func() { r.update_station_list(STATION_LIST_SEARCH) },
			help: "Search pane",
		},
		'h': KeyMapping{
			fn:   func() { r.update_station_list(STATION_LIST_HIST) },
			help: "History pane",
		},
		'f': KeyMapping{
			fn:   func() { r.update_station_list(STATION_LIST_FAV) },
			help: "Favorites pane",
		},
		'=': KeyMapping{fn: r.favorite, help: "Add to favorites"},
		'-': KeyMapping{fn: r.remove, help: "Remove from list"},
		'?': KeyMapping{fn: r.show_help, help: "Help"},
	}

	var b strings.Builder
	for _, key := range []rune{'s', 'h', 'f', '=', '-', 'p', 'q', '?'} {
		fmt.Fprintf(&b, "%c : %s\n", key, r.keymap[key].help)
	}

	r.help_modal.SetText(b.String()).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			r.grid.RemoveItem(r.help_modal)
			r.app.SetFocus(r.station_list)
		})
}

func (r *RadioUI) show_help() {
	r.grid.AddItem(r.help_modal, 3, 0, 1, 2, 0, 0, true)
	r.app.SetFocus(r.help_modal)
}
