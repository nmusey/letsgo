package core

type IModel interface {
    Table() string 
    AllColumns() string 
}

type RouterContext struct {
    App *Router
    DB  *Database
}

type Router interface {
    RegisterRoutes(ctx *RouterContext)
}
