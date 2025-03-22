package main

import (
	"fmt"

	"github.com/andrew-tawfik/outreach-routing/pkg/api"
	"github.com/andrew-tawfik/outreach-routing/pkg/app"
	"github.com/andrew-tawfik/outreach-routing/pkg/repository"
)

func main() {
	fmt.Println(api.Drive())
	fmt.Println(app.StartServer())
	fmt.Println(repository.OpenSheets())
}
