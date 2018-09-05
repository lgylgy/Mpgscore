package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, ""+
			`sheet-extractor extracts player grades from mpg googlesheet.
The result is serialized in players.json.

`)
		flag.PrintDefaults()
	}
	credentialsFile := flag.String("credentials", "", "google authorization file")
	spreadsheetID := flag.String("spreadsheet", "", "google spreadsheet identifier")
	outputDir := flag.String("output", "", "output directory")

	flag.Parse()

	b, err := ioutil.ReadFile(*credentialsFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)
	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	resp, err := srv.Spreadsheets.Get(*spreadsheetID).IncludeGridData(true).Context(context.Background()).Do()
	if err != nil {
		log.Fatal(err)
	}

	// Extract all player names
	players, err := extractPlayers(resp.Sheets)
	if err != nil {
		log.Fatal(err)
	}

	// Serialize players in players
	bytes, err := json.Marshal(players)
	if err != nil {
		log.Fatalf("Unable to marshal players: %v", err)
	}

	// Write output file
	err = ioutil.WriteFile(filepath.Join(*outputDir, "players.json"), bytes, 0644)
	if err != nil {
		log.Fatalf("Unable to write players.json: %v", err)
	}
}
