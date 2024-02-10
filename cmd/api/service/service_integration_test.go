package service_test

import (
	"go-playground/cmd/api/service"
	"go-playground/configs"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	svr    *httptest.Server
	client = &http.Client{Timeout: 1 * time.Second}
)

func TestMain(m *testing.M) {
	prop, _ := configs.Load("local")
	svc := service.New(prop, nil)
	svr = httptest.NewServer(svc)
	defer svr.Close()
	m.Run()
}

func TestNoResources(t *testing.T) {
	tests := map[string]struct {
		inMethod    string
		inResources string
		wantCode    int
		wantBody    []byte
	}{
		"no resources with get method and then 404": {
			inMethod:    "GET",
			inResources: "/",
			wantCode:    404,
			wantBody:    []byte{},
		},
		"no resources with post method and then 404": {
			inMethod:    "POST",
			inResources: "/",
			wantCode:    404,
			wantBody:    []byte{},
		},
		"no resources with put method and then 404": {
			inMethod:    "PUT",
			inResources: "/",
			wantCode:    404,
			wantBody:    []byte{},
		},
		"no resources with delete method and then 404": {
			inMethod:    "DELETE",
			inResources: "/",
			wantCode:    404,
			wantBody:    []byte{},
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			req, err := http.NewRequest(v.inMethod, testUriWith(v.inResources), nil)
			assert.NoError(t, err)
			res, err := client.Do(req)
			assert.NoError(t, err)
			got, err := io.ReadAll(res.Body)
			defer res.Body.Close()
			assert.NoError(t, err)
			assert.Equal(t, v.wantCode, res.StatusCode)
			assert.Equal(t, v.wantBody, got)
		})
	}
}

func TestStatusHandler(t *testing.T) {
	tests := map[string]struct {
		inMethod    string
		inResources string
		wantCode    int
		wantBody    []byte
	}{
		"/statuses and then success": {
			inMethod:    "GET",
			inResources: "/statuses",
			wantCode:    200,
			wantBody:    []byte("{\"status\":\"ok\"}\n"),
		},
		"/statuses with post method and then 405": {
			inMethod:    "POST",
			inResources: "/statuses",
			wantCode:    405,
			wantBody:    []byte{},
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			req, err := http.NewRequest(v.inMethod, testUriWith(v.inResources), nil)
			assert.NoError(t, err)
			res, err := client.Do(req)
			assert.NoError(t, err)
			got, err := io.ReadAll(res.Body)
			defer res.Body.Close()
			assert.NoError(t, err)
			assert.Equal(t, v.wantCode, res.StatusCode)
			assert.Equal(t, v.wantBody, got)
		})
	}
}

func TestHelloHandler_GetName(t *testing.T) {
	tests := map[string]struct {
		inMethod    string
		inResources string
		inHeader    map[string]string
		wantCode    int
		wantBody    []byte
	}{
		"/hello with valid token and then success": {
			inMethod:    "GET",
			inResources: "/hello",
			inHeader:    map[string]string{"Authorization": "admin"},
			wantCode:    200,
			wantBody:    []byte("{\"message\":\"Hello tecchu11(ADMIN)!! You have Admin role.\"}\n"),
		},
		"/hello with invalid token and then 401": {
			inMethod:    "GET",
			inResources: "/hello",
			inHeader:    map[string]string{"Authorization": "no-auth"},
			wantCode:    401,
			wantBody:    []byte("{\"type\":\"not:blank\",\"title\":\"Request With No Authentication\",\"detail\":\"Request token was not found in your request header\",\"instant\":\"/hello\",\"request_id\":\"\"}\n"),
		},
		"/hello with missing token and then 401": {
			inMethod:    "GET",
			inResources: "/hello",
			inHeader:    map[string]string{},
			wantCode:    401,
			wantBody:    []byte("{\"type\":\"not:blank\",\"title\":\"Request With No Authentication\",\"detail\":\"Request token was not found in your request header\",\"instant\":\"/hello\",\"request_id\":\"\"}\n"),
		},
		"/hello with post method and then 405": {
			inMethod:    "POST",
			inResources: "/hello",
			inHeader:    map[string]string{"Authorization": "admin"},
			wantCode:    405,
			wantBody:    []byte{},
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			req, err := http.NewRequest(v.inMethod, testUriWith(v.inResources), nil)
			assert.NoError(t, err)
			for k, v := range v.inHeader {
				req.Header.Add(k, v)
			}
			res, err := client.Do(req)
			assert.NoError(t, err)

			got, err := io.ReadAll(res.Body)
			defer res.Body.Close()
			assert.NoError(t, err)
			assert.Equal(t, v.wantCode, res.StatusCode)
			assert.Equal(t, v.wantBody, got)
		})
	}

}

func testUriWith(endpoint string) string {
	return svr.URL + endpoint
}
