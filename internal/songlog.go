/*
Keep a rolling buffer of the most recently played songs
*/
package radio

import (
	"iter"
)

const songlogBufsize = 30

type SongLog struct {
	songs [songlogBufsize]string
	idx   int
}

func (sl *SongLog) Add(song string) bool {
	iold := sl.idx - 1
	if iold < 0 {
		iold = songlogBufsize - 1
	}
	if song == sl.songs[iold] {
		return false
	}

	sl.songs[sl.idx] = song
	sl.idx++
	if sl.idx == songlogBufsize {
		sl.idx = 0
	}
	return true
}

func (sl *SongLog) Songs() iter.Seq[string] {
	return func(yield func(string) bool) {
		for i := sl.idx - 1; i >= 0; i-- {
			if sl.songs[i] == "" {
				return
			} else {
				yield(sl.songs[i])
			}
		}
		for i := songlogBufsize - 1; i >= sl.idx; i-- {
			if sl.songs[i] == "" {
				return
			} else {
				yield(sl.songs[i])
			}
		}

	}
}
