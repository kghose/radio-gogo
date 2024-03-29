package radio


type Radio struct {
	Servers        []Server
	Tags           []Tag
	Stations       []Station
	CurrentStation Station
	CurrentServer  Server
	Volume         int
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

type Tag struct {
	Name         string
	Stationcount int
}


func (r *Radio) Refresh_servers() {
	r.Servers = Get_list_of_available_servers()
}

func (r *Radio) Refresh_tags() {
	r.Tags = Get_tags(r.CurrentServer)
}


