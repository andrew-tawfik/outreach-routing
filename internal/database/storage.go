package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
)

type Database struct {
	sheet spreadsheet.Spreadsheet
}

func NewSheetClient(spreadsheetID string) (*Database, error) {

	path := os.Getenv("GOOGLE_CREDENTIALS_PATH")

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln("cannot read client_secret json file")
	}

	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	if err != nil {
		log.Fatalln("cannot config from json")
	}
	client := conf.Client(context.Background())
	service := spreadsheet.NewServiceWithClient(client)

	sheet, err := service.FetchSpreadsheet(spreadsheetID)

	if err != nil {
		log.Fatalf("Failed to fetch spreadsheet: %v", err)
	}

	fmt.Println("Connection to database established ")
	return &Database{sheet: sheet}, nil
}

func ExtractIDFromURL(url string) (string, error) {
	const marker = "/d/"

	// Locate the starting index of the spreadsheet ID
	start := strings.Index(url, marker)
	if start == -1 {
		return "", fmt.Errorf("failed to extract spreadsheet ID: '/d/' not found in URL")
	}
	start += len(marker)

	// Locate the next "/" after the ID to determine the endpoint
	end := strings.Index(url[start:], "/")
	if end == -1 {
		return "", fmt.Errorf("failed to extract spreadsheet ID: no trailing '/' found after ID")
	}

	// Extract and return the ID substring
	return url[start : start+end], nil
}
