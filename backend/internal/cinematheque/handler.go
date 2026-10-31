package cinematheque

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Movies(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/cinematheque/movies" {
		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/cinematheque/movies/"))
		if err != nil || id < 1 || strings.HasSuffix(r.URL.Path, "/") {
			writeError(w, http.StatusNotFound, "movie not found")
			return
		}
		h.movieByID(w, r, id)
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.list(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	var tags []string
	for _, tag := range strings.Split(r.URL.Query().Get("tags"), ",") {
		if tag = strings.TrimSpace(tag); tag != "" {
			tags = append(tags, tag)
		}
	}
	movies, err := h.svc.GetMovies(tags)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, movies)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var movie Movie
	if err := decodeJSON(w, r, &movie); err != nil {
		return
	}
	created, err := h.svc.CreateMovie(movie)
	if errors.Is(err, ErrValidation) {
		writeError(w, http.StatusBadRequest, "title, director, year, and rating are invalid")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) movieByID(w http.ResponseWriter, r *http.Request, id int) {
	switch r.Method {
	case http.MethodGet:
		movie, err := h.svc.GetMovie(id)
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "movie not found")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		writeJSON(w, http.StatusOK, movie)
	case http.MethodPut:
		var movie Movie
		if err := decodeJSON(w, r, &movie); err != nil {
			return
		}
		updated, err := h.svc.UpdateMovie(id, movie)
		if errors.Is(err, ErrValidation) {
			writeError(w, http.StatusBadRequest, "title, director, year, and rating are invalid")
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "movie not found")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		writeJSON(w, http.StatusOK, updated)
	case http.MethodDelete:
		err := h.svc.DeleteMovie(id)
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "movie not found")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func decodeJSON(w http.ResponseWriter, r *http.Request, destination any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request")
		return err
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
