package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"hzycoder.com/go-gin-template/internal/auth"
	"hzycoder.com/go-gin-template/internal/config"
	"hzycoder.com/go-gin-template/internal/model"
	dto "hzycoder.com/go-gin-template/internal/model/dto/request"
	"hzycoder.com/go-gin-template/internal/repository"
	"hzycoder.com/go-gin-template/pkg/response"
)

func Login(ctx context.Context, req dto.LoginUser) (string, error) {
	user, err := repository.GetUser(ctx, req.Username)
	if err != nil {
		slog.Error("query user failed", "error", err)
		return "", response.NewBizError(response.CodeUserNotFound)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)

	if err != nil {
		return "", response.NewBizError(response.CodePasswordWrong)
	}

	token, err := auth.GenerateToken(
		[]byte(config.Global.Jwt.Secret),
		user.ID,
		user.Username,
		user.Role,
	)

	if err != nil {
		slog.Error("generate token failed", "error", err)
		return "", response.NewBizError(response.CodeSystemError, "登录失败，请稍后重试")
	}

	return token, nil
}

func Register(ctx context.Context, req dto.RegisterUser) (string, error) {
	exists, err := repository.IsUserExists(ctx, req.Username)

	if err != nil {
		slog.Error("register aborted due to db error", "error", err)
		return "", response.NewBizError(response.CodeDBError, "系统繁忙，请稍后再试")
	}

	if exists {
		return "", response.NewBizError(response.CodeUserExists)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return "", response.NewBizError(response.CodeSystemError, "系统内部错误")
	}

	newUser, err := repository.AddUser(ctx, model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Role:     *req.Role,
	})

	if err != nil {
		if isDuplicateKeyError(err) {
			return "", response.NewBizError(response.CodeUserExists)
		}
		slog.Error("add user failed", "error", err)
		return "", response.NewBizError(response.CodeDBError, "注册失败，请稍后重试")
	}

	slog.Info("user registered successfully", "user_id", newUser.ID, "username", newUser.Username)

	token, err := auth.GenerateToken(
		[]byte(config.Global.Jwt.Secret),
		newUser.ID,
		newUser.Username,
		newUser.Role,
	)

	if err != nil {
		slog.Error("generate token failed", "error", err)
		return "", response.NewBizError(response.CodeSystemError, "注册成功但登录失败，请重新登录")
	}

	return token, nil
}

func isDuplicateKeyError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return true
	}

	return strings.Contains(err.Error(), "Duplicate")
}
