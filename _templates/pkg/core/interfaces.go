package core

type IModel interface {
    Table() string 
    AllColumns() string 
}

type RouterContext struct {
    DB  Database
}

type Router interface {
    RegisterRoutes()
    Serve()
}
