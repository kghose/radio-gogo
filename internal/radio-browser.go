package radio

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
)

const (
	RadioGoGoAppString     = "radiogogo"
	RadioGoGoVersion       = "0.1"
	Radio_browser_info_url = "all.api.radio-browser.info"
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

func Advanced_station_search(tag_list []string, server Server) (*StationSet, error) {
	station_set := NewStationSet()

	url := server.Name + "/json/stations/search"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return station_set, err
	}

	req.Header.Add(
		"User-Agent",
		fmt.Sprintf("%s/%s", RadioGoGoAppString, RadioGoGoVersion))
	q := req.URL.Query()
	q.Add("tagList", strings.Join(tag_list, ","))
	q.Add("hidebroken", strconv.FormatBool(true))

	req.URL.RawQuery = q.Encode()
	res, err := new(http.Client).Do(req)
	if err != nil {
		return station_set, err
	}
	defer res.Body.Close()
	var station_list []Station
	err = json.NewDecoder(res.Body).Decode(&station_list)
	if err != nil {
		return station_set, err
	}

	for i := range station_list {
		station_set.add(&station_list[i])
	}

	return station_set, nil
}
