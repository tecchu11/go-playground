package problemdetails

import (
	"encoding/json"
	"net/http"
)

type problemDetails[T any] struct {
	Type       string `json:"type"`
	Title      string `json:"title"`
	Status     int    `json:"status"`
	Detail     string `json:"detail,omitempty"`
	Instance   string `json:"instance"`
	Additional T      `json:"additional,omitempty"`
}

type Builder[T any] interface {
	WithType(string) Builder[T]
	WithDetail(string) Builder[T]
	WithAdditional(T) Builder[T]
	JSON(http.ResponseWriter, *http.Request)
}

func NewWithType[T any](title string, status int) Builder[T] {
	return &problemDetails[T]{
		Type:   "about:blank",
		Title:  title,
		Status: status,
	}
}

func New(title string, status int) Builder[any] {
	return &problemDetails[any]{
		Type:   "about:blank",
		Title:  title,
		Status: status,
	}
}

func (pb *problemDetails[T]) WithType(typ string) Builder[T] {
	pb.Type = typ
	return pb
}

func (pb *problemDetails[T]) WithDetail(detail string) Builder[T] {
	pb.Detail = detail
	return pb
}

func (pb *problemDetails[T]) WithAdditional(t T) Builder[T] {
	pb.Additional = t
	return pb
}

func (pb *problemDetails[T]) JSON(w http.ResponseWriter, r *http.Request) {
	pb.Instance = r.URL.Path
	w.Header().Add("Content-Type", "application/problem+json")
	w.WriteHeader(pb.Status)
	json.NewEncoder(w).Encode(pb)
}
