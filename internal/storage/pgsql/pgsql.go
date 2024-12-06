package pgsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"grpc-auth/internal/domain/models"
	"grpc-auth/internal/storage"
	"time"
)

type Storage struct {
	conn *pgx.Conn
}

func New(dsn string) (*Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &Storage{conn: conn}, nil
}

func (s *Storage) Close() error {
	if s.conn != nil {
		return s.conn.Close(context.Background())
	}
	return nil
}

// SaveUser сохраняет информацию о пользователе в базе данных.
func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	op := "storage.pgsql.SaveUser"
	query := `
		INSERT INTO users (email, pass_hash)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int64
	err := s.conn.QueryRow(ctx, query, email, passHash).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Код ошибки уникального ограничения
			return 0, fmt.Errorf("user already exists: %w", storage.ErrUserExists)
		}

		return 0, fmt.Errorf("failed to save user: %s, %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUser(ctx context.Context, email string) (models.User, error) {
	op := "storage.pgsql.GetUser"
	query := `
		SELECT id, email, pass_hash
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := s.conn.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s", storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("failed to get user: %s %w", op, err)
	}

	return user, nil
}

func (s *Storage) GetApp(ctx context.Context, appId int) (models.App, error) {
	op := "storage.pgsql.GetApp"
	query := `
		SELECT id, name, secret
		FROM apps
		WHERE id = $1
	`

	var app models.App
	err := s.conn.QueryRow(ctx, query, appId).Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.App{}, storage.ErrAppNotFound
		}

		return models.App{}, fmt.Errorf("failed to get app: %s %w", op, err)
	}

	return app, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	op := "storage.pgsql.IsAdmin"
	query := `
		SELECT is_admin
		FROM users
		WHERE id = $1
	`

	var isAdmin bool
	err := s.conn.QueryRow(ctx, query, userId).Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("user %d is not an admin: %w", userId, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("failed to check if user %d is an admin: %s, %w", userId, op, err)
	}

	return isAdmin, nil
}
