package service

import (
	"context"

	"hzycoder.com/go-gin-template/internal/model"
	"hzycoder.com/go-gin-template/internal/repository"
)

func GetUser(ctx context.Context, username string) (*model.User, error) {
	return repository.GetUser(ctx, username)
}
