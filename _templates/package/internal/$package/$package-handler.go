package $packages

import (
    "$appRepo/internal/core"
)

type $modelHandler struct {
    router core.Router
}

func New$modelHander(router core.Router) $modelHandler {
    return $modelHandler {
        router: router,
    }
}
