package repository

import (
	"hzycoder.com/go-gin-template/internal/database"
	"hzycoder.com/go-gin-template/internal/model"
)

func GetUser(username string) (*model.User, error) {
	row := database.DB.QueryRow(
		"SELECT id,username,password,nickname FROM users WHERE username=?",
		username,
	)

	var u model.User

	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Nickname)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
