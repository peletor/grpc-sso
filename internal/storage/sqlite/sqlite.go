package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"grpc-sso/internal/domain/models"
	"grpc-sso/internal/storage"
)

type Storage struct {
	db *sql.DB
}

// New creates a new instance of the SQLite3 storage
func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	return &Storage{db: db}, nil
}

//var _ auth.UserSaver = &Storage{}
//var _ auth.UserProvider = &Storage{}
//var _ auth.AppProvider = &Storage{}

func (s Storage) SaveUser(ctx context.Context, email string, passHash []byte) (userID int64, err error) {
	const op = "storage.sqlite.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users (email, pass_hash) VALUES (?, ?)")
	if err != nil {
		return models.EmptyUserID, fmt.Errorf("%s : %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return models.EmptyUserID, fmt.Errorf("%s : %w", op, storage.ErrUserExists)
		}

		return models.EmptyUserID, fmt.Errorf("%s : %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.EmptyUserID, fmt.Errorf("%s : %w", op, err)
	}

	return id, nil
}

// User returns user by email
func (s Storage) User(ctx context.Context, email string) (user models.User, err error) {
	const op = "storage.sqlite.User"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s : %w", op, err)
	}

	var resUser models.User
	err = stmt.QueryRowContext(ctx, email).Scan(&resUser.ID, &resUser.Email, &resUser.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s : %w", op, storage.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s : %w", op, err)
	}

	return resUser, nil
}

func (s Storage) IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error) {
	const op = "storage.sqlite.IsAdmin"

	stmt, err := s.db.Prepare("SELECT is_admin FROM users WHERE id = ?")
	if err != nil {
		return false, fmt.Errorf("%s : %w", op, err)
	}

	var resIsAdmin bool
	err = stmt.QueryRowContext(ctx, userID).Scan(&resIsAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s : %w", op, storage.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s : %w", op, err)
	}

	return resIsAdmin, nil
}

func (s Storage) App(ctx context.Context, appID int) (app models.App, err error) {
	const op = "storage.sqlite.App"

	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = ?")
	if err != nil {
		return models.App{}, fmt.Errorf("%s : %w", op, err)
	}

	var resApp models.App
	err = stmt.QueryRowContext(ctx, appID).Scan(&resApp.ID, &resApp.Name, &resApp.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s : %w", op, storage.ErrAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s : %w", op, err)
	}

	return resApp, nil
}
