/*
Package database ...
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"backend/config"
)

type SQLiteDB struct {
	db *sql.DB
}

func NewSQLiteDB(cfg *config.Config) (*SQLiteDB, error) {
	db, err := sql.Open("sqlite3", cfg.Database.Path)
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		return nil, errors.New("database connection failed")
	}

	timeoutMS := cfg.Database.TimeoutSecs * 1000
	pragmas := []string{
		"PRAGMA foreign_keys=ON;",
		fmt.Sprintf("PRAGMA busy_timeout=%d;", timeoutMS),
	}

	if cfg.Database.WALMode {
		pragmas = append(pragmas,
			"PRAGMA journal_mode=WAL;",
			"PRAGMA synchronous=NORMAL;",
			"PRAGMA cache_size=1000;",   // 1MB cache
			"PRAGMA temp_store=memory;", // Use memory for temp tables
		)
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			log.Printf("Failed to set PRAGMA settings in SQLite: %s\nFull error: %v", pragma, err)
			return nil, errors.New("failed to set SQLite PRAGMA settings")
		}
	}

	db.SetConnMaxLifetime(time.Duration(cfg.Database.TimeoutSecs) * time.Second)
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)

	if err := db.Ping(); err != nil {
		log.Printf("Failed to ping database: %v", err)
		return nil, errors.New("failed to ping database")
	}

	sqlite := &SQLiteDB{db: db}
	if err := sqlite.initSchema(); err != nil {
		return nil, err
	}

	return sqlite, nil
}

func (s *SQLiteDB) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS blog_posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL,
        tags TEXT NOT NULL,
		published_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_blog_posts_published_at ON blog_posts(published_at);
	`

	_, err := s.db.Exec(schema)
	if err != nil {
		log.Printf("Failed to initialize database with schema: %s\nFull error: %v", schema, err)
		return errors.New("failed to initialize schema")
	}

	return nil
}

func (s *SQLiteDB) Close() error {
	return s.db.Close()
}

func (s *SQLiteDB) GetDB() *sql.DB {
	return s.db
}
