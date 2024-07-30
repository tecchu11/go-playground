package httpx

import (
	"net/http"
	"strconv"
)

func QueryInt32(r *http.Request, key string) (int32, error) {
	query := r.URL.Query().Get(key)
	if query == "" {
		return 0, nil
	}
	parsed, err := strconv.ParseInt(query, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(parsed), nil
}
