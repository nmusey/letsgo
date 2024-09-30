package $packages

import (
    "$appRepo/internal/core"
)

type $packageHandler struct {
    router core.Router
}

func New$packageHander(router core.Router) $packageHandler {
    return $packageHandler {
        router: router,
    }
}
