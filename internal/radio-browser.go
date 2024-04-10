package radio

import (
	"encoding/json"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
)

const (
	RadioGoGoAppString     = "radiogogo"
	RadioGoGoVersion       = "0.0"
	Radio_browser_info_url = "all.api.radio-browser.info"
	Radio_browser_url      = "http://all.api.radio-browser.info/json/stations/search"
)

func Get_list_of_available_servers() ([]Server, error) {
	servers := []Server{}
	ips, err := net.LookupIP(Radio_browser_info_url)
	if err != nil {
		return servers, err
	}
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
	return servers, nil
}

func Pick_random_server(servers []Server) Server {
	return servers[rand.Intn(len(servers))]
}

func Advanced_station_search(tag_list []string, server Server) ([]Station, error) {
	url := server.Name + "/json/stations/search"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []Station{}, err
	}

	q := req.URL.Query()
	q.Add("tagList", strings.Join(tag_list, ","))
	q.Add("hidebroken", strconv.FormatBool(true))

	req.URL.RawQuery = q.Encode()
	res, err := new(http.Client).Do(req)
	if err != nil {
		return []Station{}, err
	}
	defer res.Body.Close()
	var t []Station
	err = json.NewDecoder(res.Body).Decode(&t)
	if err != nil {
		return []Station{}, err
	}
	return t, nil
}
