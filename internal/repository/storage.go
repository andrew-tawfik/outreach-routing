package repository

import (
	"context"
	"fmt"
	"log"
	"os"

	"strings"

	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
)

func main() {

	data, err := os.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalln("cannot read client_secret json file")
	}
	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	if err != nil {
		log.Fatalln("cannot config from json")
	}
	client := conf.Client(context.Background())
	service := spreadsheet.NewServiceWithClient(client)

	spreadsheetID := "1o90DgjjbtkFoNFkIeGyIaDpHUzlmbS0JaUbs8eTXrxA"
	sheet, err := service.FetchSpreadsheet(spreadsheetID)

	if err != nil {
		log.Fatalf("Failed to fetch spreadsheet: %v", err)
	}

	for _, row := range sheet.Sheets[0].Rows {
		for _, cell := range row {
			fmt.Printf("%s\t", cell.Value)
		}
		fmt.Println()
	}
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
