package core

import "database/sql"

type RouterContext struct {
    DB *sql.DB
}

type Router interface {
    RegisterRoutes(ctx *RouterContext)
}
