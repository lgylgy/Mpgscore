package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"mpgscore/api"
)

func createPlayer(host string, input *api.Player) error {
	output := &api.Player{}
	err := postJson(host, "/players", input, output)
	if err != nil {
		return err
	}
	fmt.Printf("%v exported...\n", output.Name)
	return nil
}

func updatePlayer(host string, input *api.Player) error {
	output := &api.Player{}
	err := putJson(host, fmt.Sprintf("/players/%s", input.ID), input, output)
	if err != nil {
		return err
	}
	fmt.Printf("%v updated...\n", output.Name)
	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, ""+
			`client post all data, contained in the input file, to a mpg server.
`)
		flag.PrintDefaults()
	}
	inputFile := flag.String("file", "", "json players file")
	host := flag.String("host", "localhost:8080", "host:port of the mpg server")
	flag.Parse()

	file, err := os.Open(*inputFile)
	if err != nil {
		log.Fatalf("Unable to open file: %v", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Unable to read file: %v", err)
	}

	var updated []*api.Player
	err = json.Unmarshal(data, &updated)
	if err != nil {
		log.Fatalf("Cannot parse players: %v", err)
	}

	result := &[]*api.Player{}
	err = getJson(*host, "/players", result)
	if err != nil {
		log.Fatalf("Cannot to fetch  players: %v", err)
	}

	actual := map[string]*api.Player{}
	for _, v := range *result {
		actual[v.Name] = v
	}

	for _, v := range updated {
		p, ok := actual[v.Name]
		if !ok {
			err := createPlayer(*host, v)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			v.ID = p.ID
			err := updatePlayer(*host, v)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
