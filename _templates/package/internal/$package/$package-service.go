package $packages

import (
    "$appRepo/internal/core"
)

type SQL$modelService struct {
    router core.Router
}

func NewSQL$modelService(router core.Router) SQL$modelService {
    return SQL$modelService {
        router: router,
    }
}
