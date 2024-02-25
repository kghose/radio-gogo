package main

import (
	"fmt"

	radio "github.com/kghose/radio-go-go/internal"
)

func main() {
	servers := radio.Get_list_of_available_servers()
	fmt.Print(servers)

	server := radio.Pick_random_server(servers)
	fmt.Print(server)

	tags := radio.Get_tags(server)
	fmt.Print(tags)

	stations := radio.Get_stations_by_tag("jazz", server)
	fmt.Print(stations)
}
