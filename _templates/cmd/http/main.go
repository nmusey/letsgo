package main

import (
    "net/http"

	"$appRepo/internal/core"
)

func main() {
    router := &core.Router{
        Mux: http.NewServeMux(),
        DB: core.NewDatabaseConnection(core.GetDefaultDatabaseConfig()),
        Cache: core.ConnectCache(core.GetDefaultCacheServer()),
    }

    BuildRoutes(router)
    router.Serve()
}
