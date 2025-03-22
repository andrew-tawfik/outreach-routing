package main

import (
	"fmt"

	"github.com/andrew-tawfik/outreach-routing/pkg/api"
	"github.com/andrew-tawfik/outreach-routing/pkg/app"
	"github.com/andrew-tawfik/outreach-routing/pkg/repository"
)

func main() {
	fmt.Println(api.StartServer())
	fmt.Println(app.Drive())
	fmt.Println(repository.OpenSheets())
}
