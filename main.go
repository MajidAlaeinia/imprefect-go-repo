package main

import (
	"github.com/vandario/govms-ipg/routes"
)

func main() {
	apiRoutes := routes.ApiRoutes()
	apiRoutes.Run(":9000")
}
