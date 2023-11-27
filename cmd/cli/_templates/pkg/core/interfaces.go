package core

type IModel interface {
    Table() string 
    Columns() string 
    ColumnValues() []interface{}
    Populate() []interface{}
}

type RouterContext struct {
    App *Router
    DB *Database
}

type Router interface {
    RegisterRoutes(ctx *RouterContext)
}
