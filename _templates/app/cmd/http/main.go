package main

import (
    "net/http"

	"$appRepo/internal/core"
    "$appRepo/cmd/http/routes"
)

func main() {
    router := &core.Router{
        Mux: http.NewServeMux(),
        DB: core.NewDatabaseConnection(core.GetDefaultDatabaseConfig()),
        Cache: core.ConnectCache(core.GetDefaultCacheServer()),
    }

    routes.BuildRoutes(router)
    router.Serve()
}
