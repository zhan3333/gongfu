package util

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"strings"
	"time"
)

func UUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func JSON[T any](v T) []byte {
	bs, _ := json.Marshal(v)
	return bs
}

func DBTimeToTimestamp(t sql.NullTime) int64 {
	if t.Valid {
		return t.Time.Unix()
	}
	return 0
}

func NullTime(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func Empty(s *string) bool {
	if s == nil || *s == "" {
		return true
	}
	return false
}
