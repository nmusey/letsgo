package middleware

import "github.com/nmusey/letsgo/server/http"

func MakeValidator(validator func (interface{}) error) func (ctx http.Ctx) error {
    return func (ctx http.Ctx) error {
        return validator(ctx)
    }
}
