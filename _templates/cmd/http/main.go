package main

import (
	"$appRepo/internal/core"
)

func main() {
    router := &core.Router{
        DB: core.NewDatabaseConnection(core.GetDefaultDatabaseConfig()),
    }

    BuildRoutes(router)
    router.Serve()
}
