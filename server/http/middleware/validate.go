package middleware

func MakeValidator(validator func (interface{}) err) func (ctx Ctx) err {
    return func (ctx Ctx) err {
        return validator(ctx))
    }
}
