package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/a-h/templ"

	"$appRepo/views/layouts"
)

func WriteJSON(w http.ResponseWriter, payload interface{}) error {
    w.Header().Set("Content-Type", "application/json")
    return json.NewEncoder(w).Encode(payload)
}

func ReadJSON(r *http.Request, payload interface{}) error {
    return json.NewDecoder(r.Body).Decode(payload)
}

func RenderTemplate(w http.ResponseWriter, components templ.Component) {
    layouts.MainLayout(components).Render(context.Background(), w)
}
