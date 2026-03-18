package cinematheque

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	svc *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) ListMovies(w http.ResponseWriter, _ *http.Request) {
	movies, err := h.svc.GetMovies()
	if err != nil {
		http.Error(w, "internal error", 500)
		return
	}
	json.NewEncoder(w).Encode(movies)
}
