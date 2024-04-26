package handler

import (
	"encoding/json"
	"fmt"
	"go-playground/pkg/problemdetails"
	"net/http"
	"strings"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// ResponseReply is response type of ReplyHandler.
type ResponseReply struct {
	Msg string `json:"message"`
}

// ReplyHandler replies.
var ReplyHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	defer newrelic.FromContext(r.Context()).StartSegment("Handler/Reply").End()
	name := strings.TrimSpace(r.PathValue("name"))
	if name == "" {
		problemdetails.New("Missing required variables", http.StatusBadRequest).
			WithDetail("Path variables name is required").
			Write(w, r)
		return
	}
	reply := ResponseReply{Msg: fmt.Sprintf("Hi %s", name)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reply)
})
