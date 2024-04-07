package main

import (
	radio "github.com/kghose/radio-go-go/internal"
)

func main() {

	r := radio.RadioUI{}
	r.Play()

	/*
		servers := radio.Get_list_of_available_servers()
		fmt.Println(servers)

		server := radio.Pick_random_server(servers)
		fmt.Println(server)

		fmt.Println("Enter tags separated by commas")
		var tag_string string
		fmt.Scanln(&tag_string)

		tag_list := strings.Split(tag_string, ",")
		stations, err := radio.Advanced_station_search(tag_list, server)
		if err != nil {
			fmt.Println(err)
		} else {
			// tags := radio.Get_tags(server)
			// fmt.Print(tags)

			// stations := radio.Get_stations_by_tag("jazz", server)
			for i, station := range stations {
				fmt.Printf("%d. %s\n", i, station.Name)
			}
		}

	*/
}
