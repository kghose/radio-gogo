package radio

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const STATION_METADATA_REFRESH_INTERVAL = time.Second

type RadioUI struct {
	device       Radio
	player       Player
	app          *tview.Application
	search_bar   *tview.InputField
	station_list *tview.List
	now_playing  *tview.TextView
	status_bar   *tview.TextView
}

func (r *RadioUI) Run() {

	r.player.Start()
	r.app = tview.NewApplication()
	r.setup_UI(r.app)
	go r.RefreshServers()
	go r.periodically_update_stream_metadata()
	if err := r.app.Run(); err != nil {
		r.app.Stop()
	}
	r.player.Quit()

}

func (r *RadioUI) setup_UI(app *tview.Application) {

	r.search_bar = tview.NewInputField().
		SetLabel("Search ").
		SetFieldWidth(80).
		SetFieldBackgroundColor(tcell.ColorGreenYellow).
		SetFieldTextColor(tcell.ColorBlack).
		SetDoneFunc(r.Search)
	r.station_list = tview.NewList().
		ShowSecondaryText(false).
		SetSelectedFunc(func(_ int, _ string, url string, _ rune) {
			r.play(url)
		})
	r.now_playing = tview.NewTextView()
	r.status_bar = tview.NewTextView()

	grid := tview.NewGrid().
		SetRows(1, -1, 4, 1).
		SetColumns(0).
		SetBorders(true).
		SetBordersColor(tcell.ColorGreenYellow)
	grid.AddItem(r.search_bar, 0, 0, 1, 1, 0, 0, true)
	grid.AddItem(r.station_list, 1, 0, 1, 1, 0, 0, false)
	grid.AddItem(r.now_playing, 2, 0, 1, 1, 0, 0, false)
	grid.AddItem(r.status_bar, 3, 0, 1, 1, 0, 0, false)
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
					len(r.device.Stations)))

			r.update_station_list(&r.device.Stations)
		})
	}()

	r.status_bar.SetText("Searching ...")

}

func (r *RadioUI) update_station_list(stations *StationSet) {
	r.station_list.Clear()
	for _, station := range *stations {
		r.station_list.AddItem(
			fmt.Sprintf("%s (%s)",
				station.Name, station.Url),
			station.Url, 0, nil)
	}
	r.station_list.SetCurrentItem(0)
	if r.station_list.GetItemCount() > 0 {
		r.app.SetFocus(r.station_list)
	}
}

func (r *RadioUI) play(url string) {
	resp := r.player.Play(url)
	r.status_bar.SetText(resp.Error)
	r.now_playing.SetText(url)
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
		}
	}
	return event
}
