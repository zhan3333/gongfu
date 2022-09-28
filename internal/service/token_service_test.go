package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToken_GetAccessToken(t *testing.T) {
	service := NewToken("20220915")
	token, err := service.GetAccessToken(1)
	assert.NoError(t, err)
	t.Logf(token)
}
