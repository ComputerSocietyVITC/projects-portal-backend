package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/models"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisRepo struct {
	Redis  *redis.Client
	Logger *zap.Logger
}

func (r *RedisRepo) SaveRefreshToken(token *models.UserRefreshToken) error {
	var key = fmt.Sprintf("refreshToken:%s", token.UserID)
	var args = redis.SetArgs{ExpireAt: time.Now().Add(time.Hour * 24 * 7)}

	if err := r.Redis.SetArgs(context.Background(), key, token.RefreshToken, args).Err(); err != nil {
		return err
	}

	return nil
}

func (r *RedisRepo) GetRefreshToken(userID string) (*models.UserRefreshToken, error) {
	var refreshToken = models.UserRefreshToken{UserID: userID}
	var key = fmt.Sprintf("refreshToken:%s", userID)

	val, err := r.Redis.Get(context.Background(), key).Result()

	if err == redis.Nil {
		val = ""
	} else if err != nil {
		return nil, err
	}

	refreshToken.RefreshToken = val

	return &refreshToken, nil
}

func (r *RedisRepo) DeleteRefreshToken(userID string) error {
	var key = fmt.Sprintf("refreshToken:%s", userID)

	if err := r.Redis.Del(context.Background(), key).Err(); err != nil {
		return err
	}
	return nil
}
