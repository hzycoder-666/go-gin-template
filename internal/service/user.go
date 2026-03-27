package service

import (
	"context"

	"hzycoder.com/go-gin-template/internal/model"
	"hzycoder.com/go-gin-template/internal/repository"
)

func GetUser(ctx context.Context, username string) (*model.User, error) {
	return repository.GetUser(ctx, username)
}

func GetUserById(ctx context.Context, userID int64) (*model.User, error) {
	return repository.GetUserById(ctx, userID)
}
