package sqlite

import (
	"context"
	"database/sql"
	"saveBot/lib/errwrap"
	"saveBot/storage"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// New creates new SQLite storage.
func New(path string) (db *Storage, err error) {
	defer func() { err = errwrap.WrapIfErr("can't connect to database", err) }()

	sqlite, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err := sqlite.Ping(); err != nil {
		return nil, err
	}
	return &Storage{db: sqlite}, nil
}

// Save saves page to storage.
func (s *Storage) Save(ctx context.Context, page *storage.Page) error {
	q := `INSERT INTO pages (url, user_name) VALUES (?, ?);`

	if _, err := s.db.ExecContext(ctx, q, page.URL, page.UserName); err != nil {
		return errwrap.Wrap("can't save page", err)
	}

	return nil
}

// PickRandom picks random page from storage.
func (s *Storage) PickRandom(ctx context.Context, username string) (*storage.Page, error) {
	q := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1;`

	var url string

	err := s.db.QueryRowContext(ctx, q, username).Scan(&url)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}

	if err != nil {
		return nil, errwrap.Wrap("can't pick random page", err)
	}

	return &storage.Page{
		URL:      url,
		UserName: username,
	}, nil
}

// Remove removes page from storage.
func (s *Storage) Remove(ctx context.Context, page *storage.Page) error {
	q := `DELETE FROM pages WHERE  url = ? AND user_name = ?;`
	if _, err := s.db.ExecContext(ctx, q, page.URL, page.UserName); err != nil {
		return errwrap.Wrap("can't remove page", err)
	}
	return nil
}

// IsExists checks if page exists in storage.
func (s *Storage) IsExists(ctx context.Context, page *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?;`

	var count int
	if err := s.db.QueryRowContext(ctx, q, page.URL, page.UserName).Scan(&count); err != nil {
		return false, errwrap.Wrap("can't check page existence", err)
	}
	return count > 0, nil
}

// Init creates the "pages" table if it doesn't exist.
func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT);`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return errwrap.Wrap("can't create table", err)
	}
	return nil
}
