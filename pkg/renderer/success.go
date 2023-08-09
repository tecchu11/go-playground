package renderer

import (
	"encoding/json"
	"net/http"
)

// Ok return 200 and passed body.
func Ok(w http.ResponseWriter, body any) {
	w.Header().Add(contentTypeKey, contentTypeValue)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(body)
}
