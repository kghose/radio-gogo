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
	Stations       *StationSet
	CurrentStation Station
	CurrentServer  Server
	User_data      UserData
}

func NewRadio() Radio {
	r := Radio{}
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
	r.Stations, err = Advanced_station_search(tag_list, r.CurrentServer)
	sort.Sort(r.Stations)
	return err

}

func (r *Radio) Refresh_servers() error {
	var err error
	r.Servers, err = Get_list_of_available_servers()
	return err
}

func (r *Radio) Now_playing(station *Station) {
	r.CurrentStation = *station
	r.User_data.Station_history.add(station)
}

const (
	USER_DATA_DIR  = "radio-gogo"
	USER_DATA_FILE = "stations.json"
)

func user_data_file() (string, error) {
	home := os.Getenv("HOME")
	data_home := os.Getenv("XDG_DATA_HOME")
	if data_home == "" {
		data_home = path.Join(home, ".local", "share")
	}
	err := os.MkdirAll(data_home, fs.FileMode(0777))
	return path.Join(data_home, USER_DATA_FILE), err
}

func (r *Radio) Save_user_data(fname string) error {
	fname, err := user_data_file()
	if err != nil {
		return err
	}
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	data_str, err := json.MarshalIndent(r.User_data, "", " ")
	_, err = f.Write(data_str)
	if err != nil {
		return err
	}

	return nil
}

func (r *Radio) Load_user_data(fname string) error {
	fname, err := user_data_file()
	if err != nil {
		return err
	}
	data_bytes, err := os.ReadFile(fname)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data_bytes, &r.User_data)
	if err != nil {
		return err
	}
	return nil
}
