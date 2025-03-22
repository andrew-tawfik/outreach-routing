package main

import (
	"fmt"
	"log"

	"github.com/andrew-tawfik/outreach-routing/internal/app"
	"github.com/andrew-tawfik/outreach-routing/internal/http"
	"github.com/andrew-tawfik/outreach-routing/internal/repository"
)

func main() {

	fmt.Println(http.StartServer())
	fmt.Println(app.Drive())
	//fmt.Println(repository.OpenSheets())

	// Step 1. Open repository and initialize program
	//Retreive guests names and addresses who will require a service
	fmt.Println("Please provide Google SheetURL")

	var googleSheetURL string
	fmt.Scanln(&googleSheetURL)

	RepositoryId, err := repository.ExtractIDFromURL(googleSheetURL)

	if err != nil {
		log.Fatalf("error extracting ID: %v", err)
	}
	fmt.Println(RepositoryId)

	// Feed Sheets url

	// Step 2. Fetch addresses exact coordinates that will be utilized

	// Step 3. Fetch distance matrix

	// Step 4. Determine the best route with RSP algorithm

	// Step 5. Output: display the routes
}
