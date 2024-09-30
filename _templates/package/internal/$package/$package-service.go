package $packages

import (
    "$appRepo/internal/core"
)

type SQL$packageService struct {
    router core.Router
}

func NewSQL$packageService(router core.Router) SQL$packageService {
    return SQL$packageService {
        router: router,
    }
}
