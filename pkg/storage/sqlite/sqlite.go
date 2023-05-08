package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"vk_chat_bot/pkg/storage"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Save(m *storage.Movie) error {
	q := `INSERT INTO movies (url, user_id) VALUES (?, ?)`

	if _, err := s.db.Exec(q, m.Url, m.UserID); err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}

	return nil
}

func (s *Storage) PickRandom(userId int) (*storage.Movie, error) {
	q := `SELECT url FROM movies WHERE user_id = ? ORDER BY RANDOM() LIMIT 1`

	var url string

	err := s.db.QueryRow(q, userId).Scan(&url)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, fmt.Errorf("can't pick random movie: %w", err)
	}

	return &storage.Movie{
		Url:    url,
		UserID: userId,
	}, nil
}

func (s *Storage) Remove(movie *storage.Movie) error {
	q := `DELETE FROM movies WHERE url = ? AND user_id = ?`
	if _, err := s.db.Exec(q, movie.Url, movie.UserID); err != nil {
		return fmt.Errorf("can't remove page: %w", err)
	}

	return nil
}

func (s *Storage) IsExist(movie *storage.Movie) (bool, error) {
	q := `SELECT COUNT(*) FROM movies WHERE url = ? AND user_id = ?`

	var count int

	if err := s.db.QueryRow(q, movie.Url, movie.UserID).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if page exists: %w", err)
	}

	return count > 0, nil
}

func (s *Storage) Init() error {
	q := `CREATE TABLE IF NOT EXISTS movies (url TEXT, user_id INTEGER)`

	_, err := s.db.Exec(q)
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	return nil
}
