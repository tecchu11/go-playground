package retrytransport_test

import (
	"errors"
	"go-playground/pkg/retrytransport"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXxx(t *testing.T) {
	tests := map[string]struct {
		inRes *http.Response
		inErr error
		want  bool
	}{
		"error is io.ErrUnexpectedEOF": {inErr: io.ErrUnexpectedEOF, want: true},
		"error is unknown":             {inErr: errors.New("this is unknown"), want: false},
	}
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			got := retrytransport.DefaultCondition(v.inRes, v.inErr)
			assert.Equal(t, v.want, got)
		})
	}
}
