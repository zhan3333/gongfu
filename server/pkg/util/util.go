package util

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"strings"
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
