package service

import (
	"hzycoder.com/go-gin-template/internal/repository"
)

func GetUser(username string) (any, error) {
	return repository.GetUser(username)
}
