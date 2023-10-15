package handler

import (
    "net/http"
)

func (h *Handler) handleError(w http.ResponseWriter, r *http.Request, err error) {
    errorMessage := os.Getenv("ENVIRONMENT") == "dev" ? err.Error() : "Internal Server Error"
    http.Error(w, errorMessage, http.StatusInternalServerError)
}
