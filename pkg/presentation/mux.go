package presentation

import (
	"net/http"
)

type MuxBuidler struct {
	mux *http.ServeMux
}

func NewMuxBuilder() *MuxBuidler {
	return &MuxBuidler{mux: &http.ServeMux{}}
}

func (builder *MuxBuidler) SetHandlerFunc(pattern string, handlerFunc http.HandlerFunc) *MuxBuidler {
	builder.mux.HandleFunc(pattern, handlerFunc)
	return builder
}

func (builder *MuxBuidler) SetHadler(pattern string, handler http.Handler) *MuxBuidler {
	builder.mux.Handle(pattern, handler)
	return builder
}

func (builder *MuxBuidler) Build() *http.ServeMux {
	return builder.mux
}
