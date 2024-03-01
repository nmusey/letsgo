package core

import (
    "net/http"
    "encoding/json"
)

func WriteJSON(w http.ResponseWriter, payload interface{}) error {
    w.Header().Set("Content-Type", "application/json")
    return json.NewEncoder(w).Encode(payload)
}

func ReadJSON(r *http.Request, payload interface{}) error {
    return json.NewDecoder(r.Body).Decode(payload)
}

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) error {
    // TODO: Implement template rendering with templ
    return nil
}
