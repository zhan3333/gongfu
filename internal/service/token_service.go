package service

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"gongfu/internal/config"
	"strings"
)

type Token interface {
	GetAccessToken(userID uint) (string, error)
	ParseAccessToken(accessToken string) (uint, error)
}

type token struct {
	Secret string
}

func NewTokenService(conf *config.Config) Token {
	return &token{Secret: conf.Token.Secret}
}

type Claims struct {
	UserID uint
	jwt.RegisteredClaims
}

func (t token) GetAccessToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID:           userID,
		RegisteredClaims: jwt.RegisteredClaims{},
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(t.Secret))

	if err != nil {
		return "", fmt.Errorf("create token: %w", err)
	}
	return tokenString, nil
}

func (t token) ParseAccessToken(accessToken string) (uint, error) {
	sp := strings.Split(accessToken, " ")
	if len(sp) != 2 {
		return 0, fmt.Errorf("access token format error")
	}
	accessToken = strings.Split(accessToken, " ")[1]
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.Secret), nil
	})
	if err != nil {
		return 0, fmt.Errorf("parse access token: %w", err)
	}
	return token.Claims.(*Claims).UserID, nil
}
