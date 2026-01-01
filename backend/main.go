package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// ============== MODELS ==============
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	IsAdmin      bool      `json:"is_admin"`
	CreatedAt    time.Time `json:"created_at"`
}

type Book struct {
	ID            string    `json:"id"`
	CoverImageURL *string   `json:"cover_image_url"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	ISBN          *string   `json:"isbn"`
	Rating        int       `json:"rating"`
	Status        string    `json:"status"`
	Tags          []string  `json:"tags"`
	CreatedAt     time.Time `json:"created_at"`
}

type Film struct {
	ID             string    `json:"id"`
	PosterImageURL *string   `json:"poster_image_url"`
	Title          string    `json:"title"`
	Rating         int       `json:"rating"`
	Tags           []string  `json:"tags"`
	CreatedAt      time.Time `json:"created_at"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// ============== DATABASE ==============
var db *sql.DB

func initDB() error {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "akinweb"),
	)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db.Ping()
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// ============== AUTH MIDDLEWARE ==============
var jwtSecret = []byte(getEnv("JWT_SECRET", "your-secret-key-change-in-production"))

func generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Next()
	}
}

// ============== AUTH HANDLERS ==============
func login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	err := db.QueryRow(
		"SELECT id, username, password_hash, is_admin, created_at FROM users WHERE username = $1",
		req.Username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.IsAdmin, &user.CreatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{Token: token, User: user})
}

// ============== BOOK HANDLERS ==============
func getBooks(c *gin.Context) {
	rows, err := db.Query(`
		SELECT id, cover_image_url, title, author, isbn, rating, created_at 
		FROM books 
		ORDER BY created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	books := []Book{}
	for rows.Next() {
		var book Book
		err := rows.Scan(
			&book.ID, &book.CoverImageURL, &book.Title,
			&book.Author, &book.ISBN, &book.Rating, &book.CreatedAt,
		)
		if err != nil {
			continue
		}

		// Get tags
		tagRows, _ := db.Query(`
			SELECT bt.name 
			FROM book_tags bt
			JOIN book_tag_links btl ON bt.id = btl.tag_id
			WHERE btl.book_id = $1
		`, book.ID)
		defer tagRows.Close()

		book.Tags = []string{}
		for tagRows.Next() {
			var tag string
			if tagRows.Scan(&tag) == nil {
				book.Tags = append(book.Tags, tag)
			}
		}

		books = append(books, book)
	}

	c.JSON(http.StatusOK, books)
}

func getBook(c *gin.Context) {
	id := c.Param("id")
	var book Book

	err := db.QueryRow(`
		SELECT id, cover_image_url, title, author, isbn, rating, created_at 
		FROM books WHERE id = $1
	`, id).Scan(
		&book.ID, &book.CoverImageURL, &book.Title,
		&book.Author, &book.ISBN, &book.Rating, &book.CreatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Get tags
	tagRows, _ := db.Query(`
		SELECT bt.name 
		FROM book_tags bt
		JOIN book_tag_links btl ON bt.id = btl.tag_id
		WHERE btl.book_id = $1
	`, book.ID)
	defer tagRows.Close()

	book.Tags = []string{}
	for tagRows.Next() {
		var tag string
		if tagRows.Scan(&tag) == nil {
			book.Tags = append(book.Tags, tag)
		}
	}

	c.JSON(http.StatusOK, book)
}

func createBook(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book.ID = uuid.New().String()
	book.CreatedAt = time.Now()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction error"})
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO books (id, cover_image_url, title, author, isbn, rating, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, book.ID, book.CoverImageURL, book.Title, book.Author, book.ISBN, book.Rating, book.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	// Handle tags
	for _, tagName := range book.Tags {
		var tagID string
		err := tx.QueryRow(`
			INSERT INTO book_tags (id, name) 
			VALUES ($1, $2) 
			ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name 
			RETURNING id
		`, uuid.New().String(), tagName).Scan(&tagID)

		if err != nil {
			continue
		}

		tx.Exec(`
			INSERT INTO book_tag_links (book_id, tag_id) 
			VALUES ($1, $2) 
			ON CONFLICT DO NOTHING
		`, book.ID, tagID)
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, book)
}

func updateBook(c *gin.Context) {
	id := c.Param("id")
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction error"})
		return
	}
	defer tx.Rollback()

	result, err := tx.Exec(`
		UPDATE books 
		SET cover_image_url = $1, title = $2, author = $3, isbn = $4, rating = $5
		WHERE id = $6
	`, book.CoverImageURL, book.Title, book.Author, book.ISBN, book.Rating, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	// Update tags - delete old links and create new ones
	tx.Exec("DELETE FROM book_tag_links WHERE book_id = $1", id)

	for _, tagName := range book.Tags {
		var tagID string
		err := tx.QueryRow(`
			INSERT INTO book_tags (id, name) 
			VALUES ($1, $2) 
			ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name 
			RETURNING id
		`, uuid.New().String(), tagName).Scan(&tagID)

		if err != nil {
			continue
		}

		tx.Exec(`
			INSERT INTO book_tag_links (book_id, tag_id) 
			VALUES ($1, $2) 
			ON CONFLICT DO NOTHING
		`, id, tagID)
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	book.ID = id
	c.JSON(http.StatusOK, book)
}

func deleteBook(c *gin.Context) {
	id := c.Param("id")

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction error"})
		return
	}
	defer tx.Rollback()

	// Delete tag links first
	tx.Exec("DELETE FROM book_tag_links WHERE book_id = $1", id)

	result, err := tx.Exec("DELETE FROM books WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

// ============== FILM HANDLERS (Ready for future use) ==============
func getFilms(c *gin.Context) {
	rows, err := db.Query(`
		SELECT id, poster_image_url, title, rating, created_at 
		FROM films 
		ORDER BY created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	films := []Film{}
	for rows.Next() {
		var film Film
		err := rows.Scan(&film.ID, &film.PosterImageURL, &film.Title, &film.Rating, &film.CreatedAt)
		if err != nil {
			continue
		}

		// Get tags
		tagRows, _ := db.Query(`
			SELECT ft.name 
			FROM film_tags ft
			JOIN film_tag_links ftl ON ft.id = ftl.tag_id
			WHERE ftl.film_id = $1
		`, film.ID)
		defer tagRows.Close()

		film.Tags = []string{}
		for tagRows.Next() {
			var tag string
			if tagRows.Scan(&tag) == nil {
				film.Tags = append(film.Tags, tag)
			}
		}

		films = append(films, film)
	}

	c.JSON(http.StatusOK, films)
}

func createFilm(c *gin.Context) {
	var film Film
	if err := c.ShouldBindJSON(&film); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	film.ID = uuid.New().String()
	film.CreatedAt = time.Now()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction error"})
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO films (id, poster_image_url, title, rating, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, film.ID, film.PosterImageURL, film.Title, film.Rating, film.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create film"})
		return
	}

	// Handle tags
	for _, tagName := range film.Tags {
		var tagID string
		err := tx.QueryRow(`
			INSERT INTO film_tags (id, name) 
			VALUES ($1, $2) 
			ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name 
			RETURNING id
		`, uuid.New().String(), tagName).Scan(&tagID)

		if err != nil {
			continue
		}

		tx.Exec(`
			INSERT INTO film_tag_links (film_id, tag_id) 
			VALUES ($1, $2) 
			ON CONFLICT DO NOTHING
		`, film.ID, tagID)
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, film)
}

func deleteFilm(c *gin.Context) {
	id := c.Param("id")

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction error"})
		return
	}
	defer tx.Rollback()

	tx.Exec("DELETE FROM film_tag_links WHERE film_id = $1", id)

	result, err := tx.Exec("DELETE FROM films WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete film"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Film not found"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Film deleted successfully"})
}

// ============== MAIN ==============
func main() {
	// Initialize database
	if err := initDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Setup Gin
	router := gin.Default()

	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Public routes
	router.POST("/api/auth/login", login)
	router.GET("/api/books", getBooks)
	router.GET("/api/books/:id", getBook)
	router.GET("/api/films", getFilms)

	// Protected routes
	protected := router.Group("/api")
	protected.Use(authMiddleware())
	{
		protected.POST("/books", createBook)
		protected.PUT("/books/:id", updateBook)
		protected.DELETE("/books/:id", deleteBook)

		protected.POST("/films", createFilm)
		protected.DELETE("/films/:id", deleteFilm)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s...", port)
	router.Run(":" + port)
}
