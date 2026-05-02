/*
Functions to find and query radio browser servers
*/

package radio_browser

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"maps"
	"math/rand"
	"net"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

const (
	appString          = "radiogogo"
	appVersion         = "0.2"
	radioBrowserAPIUrl = "all.api.radio-browser.info"
	searchLimit        = "10000"
)

func GetAvailableServers() ([]string, error) {
	ips, err := net.LookupIP(radioBrowserAPIUrl)
	if err != nil {
		slog.Error("DNS lookup failed", "Host", radioBrowserAPIUrl)
		return []string{}, err
	}

	servers_set := make(map[string]struct{})
	for _, ip := range ips {
		addrs, err := net.LookupAddr(ip.String())
		if err != nil {
			slog.Error("Address look up failed", "IP", ip)
			continue
		}
		for _, addr := range addrs {
			url := "https://" + strings.TrimSuffix(addr, ".")
			servers_set[url] = struct{}{}
		}
	}
	servers := slices.Collect(maps.Keys(servers_set))
	slog.Info("Lookup servers", "found", len(servers))

	return servers, nil
}

func PickRandomServer(servers []string) string {
	server := servers[rand.Intn(len(servers))]
	slog.Info("Pick server", "url", server)
	return server
}

type Station struct {
	ChangeUUID  string  `json:"changeuuid"`
	StationUUID string  `json:"stationuuid"`
	Name        string  `json:"name"`
	URL         string  `json:"url"`
	URLResolved string  `json:"url_resolved"`
	Homepage    string  `json:"homepage"`
	Favicon     string  `json:"favicon"`
	Tags        string  `json:"tags"`
	CountryCode string  `json:"countrycode"`
	State       string  `json:"state"`
	Language    string  `json:"language"`
	Latitude    float64 `json:"geo_lat"`
	Longitude   float64 `json:"geo_long"`
}


func sanitizeStrings(stations []Station) {
	for _, station := range stations {
		station.Name = fmt.Sprintf("%q", station.Name)
		station.URLResolved = fmt.Sprintf("%q", station.URLResolved) // Maybe redundant?
		// TODO: Other strings that we may print
	}
}


func StationSearch(comma_separated_keywords string, server_url string) ([]Station, error) {

	stations := []Station{}

	url := server_url + "/json/stations/search"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Error("Setting up request failed", "Error", err)
		return stations, err
	}

	req.Header.Add(
		"User-Agent",
		fmt.Sprintf("%s/%s", appString, appVersion))
	q := req.URL.Query()
	q.Add("tagList", comma_separated_keywords)
	q.Add("hidebroken", strconv.FormatBool(true))
	q.Add("limit", searchLimit)

	req.URL.RawQuery = q.Encode()
	res, err := new(http.Client).Do(req)
	if err != nil {
		slog.Error("Search failed", "Error", err)
		return stations, err
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&stations)
	if err != nil {
		slog.Error("Error parsing station data", "Error", err)
		return stations, err
	}

	dedupedStations := dedupeStationList(stations)
	sanitizeStrings(dedupedStations)
	slog.Info(
		"Search",
		"keywords", comma_separated_keywords,
		"found", len(dedupedStations))
	return dedupedStations, nil
}

func dedupeStationList(stations []Station) []Station {
	seen := make(map[string]bool)
	deduped_stations := []Station{}
	for _, station := range stations {
		if !seen[station.URLResolved] {
			seen[station.URLResolved] = true
			deduped_stations = append(deduped_stations, station)
		}
	}
	return deduped_stations
}
