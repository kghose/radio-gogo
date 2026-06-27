/*
Keep a rolling buffer of the most recently played songs
*/
package radio

type SongLog struct {
	songs []string
	idx   int
}

func NewSongLog(n int) SongLog {
	songlog := SongLog{}
	songlog.songs = make([]string, n)
	return songlog
}

func (sl *SongLog) Add(song string) bool {
	iold := sl.idx - 1
	if iold < 0 {
		iold = len(sl.songs) - 1
	}
	if song == sl.songs[iold] {
		return false
	}

	sl.songs[sl.idx] = song
	sl.idx++
	if sl.idx == len(sl.songs) {
		sl.idx = 0
	}
	return true
}

func (sl *SongLog) Songs() []string {
	songs := []string{}
	for i := sl.idx - 1; i >= 0; i-- {
		if sl.songs[i] == "" {
			return songs
		} else {
			songs = append(songs, sl.songs[i])
		}
	}
	for i := len(sl.songs) - 1; i >= sl.idx; i-- {
		if sl.songs[i] == "" {
			return songs
		} else {
			songs = append(songs, sl.songs[i])
		}
	}
	return songs
}
