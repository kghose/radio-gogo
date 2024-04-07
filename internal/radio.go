/*
 */
package radio

import "fmt"

type EventKind string

const (
	ERROR           EventKind = "ERROR"
	DATA_REFRESHED  EventKind = "DATA_REFRESHED"
	STATE_REFRESHED EventKind = "STATE_REFRESHED"
)

type Event struct {
	kind    EventKind
	message string
}

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

func (r *Radio) Init(radio_q chan Event) {
	r.Servers = Get_list_of_available_servers()
	radio_q <- Event{
		STATE_REFRESHED,
		fmt.Sprintf("%d Radio Browser servers found.", len(r.Servers))}
}

func (r *Radio) Refresh_servers() {
	r.Servers = Get_list_of_available_servers()
}

func (r *Radio) Refresh_tags() {
	r.Tags = Get_tags(r.CurrentServer)
}
