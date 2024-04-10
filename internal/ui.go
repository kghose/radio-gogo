package radio

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type State struct {
	tag_string    string
	selected_tags []string
	station_index int
}

type RadioUI struct {
	device       Radio
	player       Player
	app          *tview.Application
	search_bar   *tview.InputField
	station_list *tview.List
	now_playing  *tview.TextView
	status_bar   *tview.TextView
	state        State
	radio_on     bool
	error_msg    string
	msg          string
}

func (r *RadioUI) Run() {

	r.player.Start()
	r.app = tview.NewApplication()
	r.setup_UI(r.app)
	go r.RefreshServers()
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
		SetRows(1, -3, -1, 1).
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
		r.device.FindByTag([]string{r.search_bar.GetText()})
		r.app.QueueUpdateDraw(func() {

			if err != nil {
				r.status_bar.SetText(
					fmt.Sprintf("Search error: %s", err.Error()))
				return
			}

			r.station_list.Clear()
			for _, station := range r.device.Stations {
				r.station_list.AddItem(
					fmt.Sprintf("%s (%s)",
						station.Name, station.Url),
					station.Url, 0, nil)
			}
			r.status_bar.SetText(
				fmt.Sprintf("Found %d stations.",
					len(r.device.Stations)))
			if len(r.device.Stations) > 0 {
				r.app.SetFocus(r.station_list)
			}
		})
	}()

	r.status_bar.SetText("Searching ...")

}

func (r *RadioUI) play(url string) {
	resp := r.player.Play(url)
	r.status_bar.SetText(resp.Error)
	r.now_playing.SetText(url)
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
