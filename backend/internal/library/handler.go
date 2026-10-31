package library

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

func (h *Handler) Books(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/library/books" {
		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/library/books/"))
		if err != nil || id < 1 || strings.HasSuffix(r.URL.Path, "/") {
			writeError(w, http.StatusNotFound, "book not found")
			return
		}
		h.bookByID(w, r, id)
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
	books, err := h.svc.GetBooks(tags)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, books)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := decodeJSON(w, r, &book); err != nil {
		return
	}
	created, err := h.svc.CreateBook(book)
	if errors.Is(err, ErrValidation) {
		writeError(w, http.StatusBadRequest, "title, author, year, and rating are invalid")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) bookByID(w http.ResponseWriter, r *http.Request, id int) {
	switch r.Method {
	case http.MethodGet:
		book, err := h.svc.GetBook(id)
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "book not found")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		writeJSON(w, http.StatusOK, book)
	case http.MethodPut:
		var book Book
		if err := decodeJSON(w, r, &book); err != nil {
			return
		}
		updated, err := h.svc.UpdateBook(id, book)
		if errors.Is(err, ErrValidation) {
			writeError(w, http.StatusBadRequest, "title, author, year, and rating are invalid")
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "book not found")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		writeJSON(w, http.StatusOK, updated)
	case http.MethodDelete:
		err := h.svc.DeleteBook(id)
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "book not found")
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
