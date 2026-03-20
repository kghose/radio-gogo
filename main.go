package main

import (
	"fmt"
	"log"
	// radio "github.com/kghose/radio-go-go/internal"
	radio_browser "github.com/kghose/radio-go-go/internal/radio_browser"
)

func main() {

	servers, err := radio_browser.GetAvailableServers()
	if err != nil {
		log.Fatal("Could not find radio browser servers")
	}

	server := radio_browser.PickRandomServer(servers)
	stations, err := radio_browser.StationSearch("jazz", server)

	for _, station := range stations {
		fmt.Println(station.Name + ":" + station.Tags)
	}
	fmt.Println(len(stations))

	// r := radio.RadioUI{}
	// r.Run()

}
