package radio

import (
	"encoding/json"
	"math/rand"
	"net"
	"net/http"
	"strings"
)

const (
	RadioGoGoAppString     = "radiogogo"
	RadioGoGoVersion       = "0.0"
	Radio_browser_info_url = "all.api.radio-browser.info"
	Radio_browser_url      = "http://all.api.radio-browser.info/json/stations/search"
)

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


func Get_list_of_available_servers() []Server {
	ips, err := net.LookupIP(Radio_browser_info_url)
	if err != nil {
		panic(err)
	}
	servers := []Server{}
	for _, ip := range ips {
		names, err := net.LookupAddr(ip.String())
		for _, name := range names {
			servers = append(servers, Server{
				Name: "https://" + strings.TrimSuffix(name, "."),
				IP:   ip.String(),
				Err:  err,
			})
		}
	}
	return servers
}

func Pick_random_server(servers []Server) Server {
	return servers[rand.Intn(len(servers))]
}

type map_type map[string]interface{}

func (t *map_type) f(s string) map_type {
	return (*t)[s].(map[string]interface{})
}

func (t *map_type) str(s string) string {
	return (*t)[s].(string)
}
func (t *map_type) f64(s string) float64 {
	return (*t)[s].(float64)
}

func Get_query[T Station | Tag](url string) []T {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	q := req.URL.Query()
	// q.Set("action", "wbgetentities")
	// q.Set("format", "json")
	// q.Set("ids", "Q24871")
	req.URL.RawQuery = q.Encode()
	res, err := new(http.Client).Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	var t []T
	err = json.NewDecoder(res.Body).Decode(&t)
	if err != nil {
		panic(err)
	}
	// json.Unmarshal(res.Body, &t)
	return t
}

func Get_tags(server Server) []Tag {
	url := server.Name + "/json/tags"
	res := Get_query[Tag](url)
	return res
}

func Get_stations_by_tag(tag string, server Server) []Station {
	url := server.Name + "/json/stations/bytag/" + tag
	res := Get_query[Station](url)
	return res
}
