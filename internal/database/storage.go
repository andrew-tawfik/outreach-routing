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

// Database wraps the Google Sheets spreadsheet object
// and represents your external data source.
type Database struct {
	sheet spreadsheet.Spreadsheet
}

// NewSheetClient initializes a Google Sheets client using credentials
func NewSheetClient(spreadsheetID string) (*Database, error) {

	path := os.Getenv("GOOGLE_CREDENTIALS_PATH")

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln("cannot read client_secret json file")
	}

	// Create a JWT-based authenticated client using the provided credentials
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

// ExtractIDFromURL parses a standard Google Sheets URL
func ExtractIDFromURL(url string) (string, error) {
	const marker = "/d/"

	start := strings.Index(url, marker)
	if start == -1 {
		return "", fmt.Errorf("failed to extract spreadsheet ID: '/d/' not found in URL")
	}
	start += len(marker)

	end := strings.Index(url[start:], "/")
	if end == -1 {
		return "", fmt.Errorf("failed to extract spreadsheet ID: no trailing '/' found after ID")
	}

	return url[start : start+end], nil
}
