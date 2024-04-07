package radio

import "github.com/nsf/termbox-go"

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
