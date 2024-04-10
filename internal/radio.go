/*
 */
package radio

type EventKind string

type Radio struct {
	Servers        []Server
	Stations       []Station
	CurrentStation Station
	CurrentServer  Server
	Error          string
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
	return err

}

func (r *Radio) Refresh_servers() error {
	var err error
	r.Servers, err = Get_list_of_available_servers()
	return err
}
