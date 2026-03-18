package internal

import (
	"database/sql"
	"net/http"

	"ak-in-web/internal/cinematheque"
	"ak-in-web/internal/library"
)

func Router(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Library
	libRepo := library.NewRepository(db)
	libSvc := library.NewService(libRepo)
	libHandler := library.NewHandler(libSvc)

	mux.HandleFunc("/api/library/books", libHandler.ListBooks)

	// Cinematheque
	cinRepo := cinematheque.NewRepository(db)
	cinSvc := cinematheque.NewService(cinRepo)
	cinHandler := cinematheque.NewHandler(cinSvc)

	mux.HandleFunc("/api/cinematheque/movies", cinHandler.ListMovies)

	return mux
}
