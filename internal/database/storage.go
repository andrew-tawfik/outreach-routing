package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/andrew-tawfik/outreach-routing/internal/config"
	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
)



type Database struct {
	sheet spreadsheet.Spreadsheet
}


func NewSheetClient(spreadsheetID string) (*Database, error) {
	embeddedData := config.GetEmbeddedServiceAccountJSON()
	if len(embeddedData) > 0 {
		conf, err := google.JWTConfigFromJSON(embeddedData, spreadsheet.Scope)
		if err == nil {
			client := conf.Client(context.Background())
			service := spreadsheet.NewServiceWithClient(client)

			sheet, err := service.FetchSpreadsheet(spreadsheetID)
			if err != nil {
			} else {
				return &Database{sheet: sheet}, nil
			}
		} else {
			fmt.Printf("Warning: embedded credentials are invalid: %v\n", err)
		}
	}

	
	projectRoot, err := filepath.Abs(filepath.Join(".", ".."))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve project root:", err)
	}
	credentialsPath := filepath.Join(projectRoot, "client_secret.json")

	data, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %v", credentialsPath, err)
	}

	
	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	if err != nil {
		return nil, fmt.Errorf("cannot config from json")
	}
	client := conf.Client(context.Background())
	service := spreadsheet.NewServiceWithClient(client)

	sheet, err := service.FetchSpreadsheet(spreadsheetID)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch spreadsheet: %v", err)
	}

	return &Database{sheet: sheet}, nil
}


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
