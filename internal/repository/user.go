package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"hzycoder.com/go-gin-template/internal/database"
	"hzycoder.com/go-gin-template/internal/model"
)

func GetUser(c context.Context, username string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	row := database.DB.QueryRowContext(
		ctx,
		"SELECT id,username,password,nickname,role FROM users WHERE username=?",
		username,
	)

	var u model.User

	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Nickname, &u.Role)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func IsUserExists(c context.Context, username string) (bool, error) {
	var dummy int
	err := database.DB.QueryRowContext(
		c,
		"SELECT 1 FROM users WHERE username = ? LIMIT 1",
		username,
	).Scan(&dummy)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		slog.Error("check user exists failed (system error)", "username", username, "error", err)
		return false, fmt.Errorf("database check failed: %w", err)
	}

	return true, nil
}

func AddUser(c context.Context, user model.User) (*model.User, error) {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	query := `
		INSERT INTO users (username, password, role, created_at, updated_at)
		VALUES (?,?,?,?,?)
	`

	result, err := database.DB.ExecContext(c, query,
		user.Username,
		user.Password,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		slog.Error("insert user failed", "error", err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err == nil {
		user.ID = id
	}

	return &user, nil
}
