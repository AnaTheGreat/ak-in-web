package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func New(host, port, name, user, password, sslMode string) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, name, sslMode,
	)
	return sql.Open("postgres", dsn)
}

func WaitForConnection(database *sql.DB, attempts int, delay time.Duration) error {
	var err error
	for i := 0; i < attempts; i++ {
		if err = database.Ping(); err == nil {
			return nil
		}
		time.Sleep(delay)
	}
	return fmt.Errorf("database connection failed after %d attempts: %w", attempts, err)
}

func RunMigrations(database *sql.DB, directory string) error {
	if _, err := database.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version BIGINT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`); err != nil {
		return fmt.Errorf("create migration table: %w", err)
	}

	entries, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	for _, file := range files {
		version, err := migrationVersion(file)
		if err != nil {
			return err
		}

		var applied bool
		if err := database.QueryRow(
			"SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version = $1)",
			version,
		).Scan(&applied); err != nil {
			return fmt.Errorf("check migration %s: %w", file, err)
		}
		if applied {
			continue
		}

		sqlBytes, err := os.ReadFile(filepath.Join(directory, file))
		if err != nil {
			return fmt.Errorf("read migration %s: %w", file, err)
		}

		tx, err := database.Begin()
		if err != nil {
			return fmt.Errorf("begin migration %s: %w", file, err)
		}
		if _, err := tx.Exec(string(sqlBytes)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("apply migration %s: %w", file, err)
		}
		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record migration %s: %w", file, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", file, err)
		}
	}

	return nil
}

func migrationVersion(filename string) (int64, error) {
	prefix := strings.SplitN(filename, "_", 2)[0]
	version, err := strconv.ParseInt(prefix, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid migration filename %q: %w", filename, err)
	}
	return version, nil
}
