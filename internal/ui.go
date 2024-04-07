package radio

import (
	"github.com/nsf/termbox-go"
)

type State struct {
	tag_string    string
	selected_tags []string
	station_index int
}

type RadioUI struct {
	state    State
	radio_on bool
}

func (r *RadioUI) Play() {

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	event_q := make(chan termbox.Event)
	go r.event_poll_loop(event_q)
	r.main_loop(event_q)
}

func (r *RadioUI) event_poll_loop(event_q chan termbox.Event) {
	for {
		event_q <- termbox.PollEvent()
	}
}

func (r *RadioUI) main_loop(event_q chan termbox.Event) {
	r.radio_on = true
	for r.radio_on {
		r.render()
		select {
		case ev := <-event_q:
			r.process_termbox_event(ev)
		}
	}
}

func (r *RadioUI) render() {
	termbox.Clear(termbox.ColorBlue, termbox.ColorBlue)
	w, _ := termbox.Size()
	banner(0, 0, w, "Radio Go Go", true, termbox.ColorBlack, termbox.ColorCyan)
	err := termbox.Flush()
	if err != nil {
		panic(err)
	}
}

func banner(
	row int,
	col int,
	width int,
	msg string,
	center bool,
	fg termbox.Attribute,
	bg termbox.Attribute,
) {
	if center {
		col = max(0, width/2-len(msg)/2)
	}
	for n := 0; n < col; n++ {
		termbox.SetCell(n, row, ' ', bg, bg)
	}
	for n, c := range msg {
		termbox.SetCell(col+n, row, c, fg, bg)
	}
	for n := col + len(msg); n < width; n++ {
		termbox.SetCell(n, row, ' ', bg, bg)
	}

}

func (r *RadioUI) process_termbox_event(ev termbox.Event) {
	if ev.Type != termbox.EventKey {
		return
	}
	if quit_key(ev) {
		r.radio_on = false
		return
	}

}

func quit_key(ev termbox.Event) bool {
	return ev.Ch == 'q' ||
		ev.Key == termbox.KeyEsc ||
		ev.Key == termbox.KeyCtrlC ||
		ev.Key == termbox.KeyCtrlD
}
