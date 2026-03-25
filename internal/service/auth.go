package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"hzycoder.com/go-gin-template/internal/auth"
	"hzycoder.com/go-gin-template/internal/config"
	"hzycoder.com/go-gin-template/internal/repository"
)

func Login(username, password string) (string, error) {
	user, err := repository.GetUser(username)
	if err != nil {
		return "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)

	if err != nil {
		return "", errors.New("password incorrect")
	}

	token, err := auth.GenerateToken(
		[]byte(config.Global.Jwt.Secret),
		user.ID,
		user.Username,
	)

	if err != nil {
		return "", err
	}

	return token, nil
}
