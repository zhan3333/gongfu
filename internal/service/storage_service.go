package service

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	_ "github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"time"
)

type Storage interface {
	GetPresignedURL(ctx context.Context, key string) (string, error)
	KeyExists(ctx context.Context, key string) (bool, error)
	GetVisitURL(ctx context.Context, key string) (string, error)
	GetPublicVisitURL(ctx context.Context, key string) string
}

type storage struct {
	Cos       *cos.Client
	SecretID  string
	SecretKey string
}

func (s storage) GetVisitURL(ctx context.Context, key string) (string, error) {
	u, err := s.Cos.Object.GetPresignedURL(ctx, http.MethodGet, key, s.SecretID, s.SecretKey, time.Hour, nil)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func (s storage) GetPublicVisitURL(ctx context.Context, key string) string {
	u := s.Cos.Object.GetObjectURL(key)
	return u.String()
}

func (s storage) KeyExists(ctx context.Context, key string) (bool, error) {
	return s.Cos.Object.IsExist(ctx, key)
}

func (s storage) GetPresignedURL(ctx context.Context, key string) (string, error) {
	u, err := s.Cos.Object.GetPresignedURL(ctx, http.MethodPut, key, s.SecretID, s.SecretKey, time.Hour, nil)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func NewStorage(cos *cos.Client, secretID string, secreteKey string) Storage {
	return &storage{Cos: cos, SecretID: secretID, SecretKey: secreteKey}
}
