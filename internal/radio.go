/*
 */
package radio

import (
	"encoding/json"
	"io/fs"
	"os"
	"path"
	"sort"
)

type StationSet struct {
	Stations []*Station
	urls     map[string]*Station
}

func (s *StationSet) add(station *Station) {
	_, present := s.urls[station.Url]
	if present {
		return
	} else {
		s.Stations = append(s.Stations, station)
		s.urls[station.Url] = station
	}
}

func (s *StationSet) By_url(url string) *Station {
	return s.urls[url]
}

// https://pkg.go.dev/sort
func (s *StationSet) Len() int { return len(s.Stations) }
func (s *StationSet) Swap(i, j int) {
	s.Stations[i], s.Stations[j] = s.Stations[j], s.Stations[i]
}
func (s *StationSet) Less(i, j int) bool {
	return s.Stations[i].Name < s.Stations[j].Name
}

func NewStationSet() *StationSet {
	s := StationSet{}
	s.urls = make(map[string]*Station)
	return &s
}

type Radio struct {
	Servers        []Server
	StationLists   [3]*StationSet
	Stations       *StationSet
	CurrentStation Station
	CurrentServer  Server
	User_data      UserData
}

const (
	STATION_LIST_SEARCH = 0
	STATION_LIST_HIST   = 1
	STATION_LIST_FAV    = 2
)

var uSER_DATA_FILE = [3]string{"", "history.json", "favorites.json"}

func NewRadio() Radio {
	r := Radio{}
	for i := 0; i < 3; i++ {
		r.StationLists[i] = NewStationSet()
	}
	r.Stations = NewStationSet()
	r.User_data = NewUserData()
	return r
}

type UserData struct {
	Station_history   *StationSet
	Station_favorites *StationSet
}

func NewUserData() UserData {
	ud := UserData{}
	ud.Station_history = NewStationSet()
	ud.Station_favorites = NewStationSet()
	return ud
}

type Server struct {
	Name string
	IP   string
	Err  error
}

type Station struct {
	Name string
	Url  string
}

func (r *Radio) FindByTag(tag_list []string) error {

	r.CurrentServer = Pick_random_server(r.Servers)

	var err error
	r.StationLists[STATION_LIST_SEARCH], err = Advanced_station_search(tag_list, r.CurrentServer)
	sort.Sort(r.StationLists[STATION_LIST_SEARCH])
	r.Stations = r.StationLists[STATION_LIST_SEARCH]
	return err

}

func (r *Radio) Refresh_servers() error {
	var err error
	r.Servers, err = Get_list_of_available_servers()
	return err
}

func (r *Radio) Now_playing(station *Station) {
	r.CurrentStation = *station
	r.StationLists[STATION_LIST_HIST].add(station)
	r.User_data.Station_history.add(station)
}

const (
	USER_DATA_DIR                    = "radio-gogo"
	USER_DATA_STATION_HISTORY_FILE   = "history.json"
	USER_DATA_STATION_FAVORITES_FILE = "favorites.json"
)

func user_data_file() (string, string, error) {
	home := os.Getenv("HOME")
	data_home := os.Getenv("XDG_DATA_HOME")
	if data_home == "" {
		data_home = path.Join(home, ".local", "share")
	}
	data_home = path.Join(data_home, USER_DATA_DIR)
	err := os.MkdirAll(data_home, fs.FileMode(0777))
	return path.Join(
			data_home,
			USER_DATA_STATION_HISTORY_FILE,
		),
		path.Join(
			data_home,
			USER_DATA_STATION_FAVORITES_FILE,
		), err
}

func get_user_data_dir() (string, error) {
	home := os.Getenv("HOME")
	data_home := os.Getenv("XDG_DATA_HOME")
	if data_home == "" {
		data_home = path.Join(home, ".local", "share")
	}
	data_home = path.Join(data_home, USER_DATA_DIR)
	err := os.MkdirAll(data_home, fs.FileMode(0777))
	return data_home, err
}

func (r *Radio) Save_user_data() error {
	user_data_dir, err := get_user_data_dir()
	if err != nil {
		return err
	}
	for _, i := range []int{STATION_LIST_FAV, STATION_LIST_HIST} {
		if err = save_user_data_file(path.Join(user_data_dir, uSER_DATA_FILE[i]), r.StationLists[i]); err != nil {
			return err
		}
	}
	return nil
}

func save_user_data_file(fname string, data *StationSet) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	data_str, err := json.MarshalIndent(data.Stations, "", " ")
	_, err = f.Write(data_str)
	if err != nil {
		return err
	}

	return nil
}

func (r *Radio) Load_user_data() error {
	user_data_dir, err := get_user_data_dir()
	if err != nil {
		return err
	}

	for _, i := range []int{STATION_LIST_FAV, STATION_LIST_HIST} {
		if err = load_user_data_file(path.Join(user_data_dir, uSER_DATA_FILE[i]), r.StationLists[i]); err != nil {
			return err
		}
	}
	return nil
}

func load_user_data_file(fname string, station_set *StationSet) error {

	data_bytes, err := os.ReadFile(fname)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	stations := []Station{}
	if err = json.Unmarshal(data_bytes, &stations); err != nil {
		return err
	}

	for i := range stations {
		station_set.add(&stations[i])
	}
	return nil
}
