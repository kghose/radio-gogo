package radio

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

type State struct {
	tag_string    string
	selected_tags []string
	station_index int
}

type RadioUI struct {
	device    Radio
	state     State
	radio_on  bool
	error_msg string
	msg       string
}

func (r *RadioUI) Play() {

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	event_q := make(chan termbox.Event)
	radio_q := make(chan Event)
	go r.event_poll_loop(event_q)
	r.main_loop(event_q, radio_q)
}

func (r *RadioUI) event_poll_loop(event_q chan termbox.Event) {
	for {
		event_q <- termbox.PollEvent()
	}
}

func (r *RadioUI) main_loop(event_q chan termbox.Event, radio_q chan Event) {
	r.radio_on = true
	go r.device.Init(radio_q)
	for r.radio_on {
		r.render()
		select {
		case ev := <-event_q:
			r.process_termbox_event(ev)

		case rev := <-radio_q:
			r.process_event(rev)
		}
	}
}

func (r *RadioUI) render() {
	termbox.Clear(termbox.ColorBlue, termbox.ColorBlue)
	w, h := termbox.Size()
	banner(0, 0, w, "Radio Go Go", true, termbox.ColorBlack, termbox.ColorCyan)

	if len(r.error_msg) > 0 {
		banner(
			h-1,
			0,
			w,
			fmt.Sprintf("Error: %s", r.error_msg),
			true,
			termbox.ColorRed,
			termbox.ColorYellow,
		)
	}
	if len(r.msg) > 0 {
		banner(
			h-1,
			0,
			w,
			r.msg,
			true,
			termbox.ColorBlack,
			termbox.ColorCyan,
		)
	}
	err := termbox.Flush()
	if err != nil {
		panic(err)
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

func (r *RadioUI) process_event(rev Event) {
	switch rev.kind {
	case ERROR:
		r.error_msg = rev.message
	case STATE_REFRESHED:
		r.msg = rev.message

	}

}

func quit_key(ev termbox.Event) bool {
	return ev.Ch == 'q' ||
		ev.Key == termbox.KeyEsc ||
		ev.Key == termbox.KeyCtrlC ||
		ev.Key == termbox.KeyCtrlD
}
