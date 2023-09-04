package service_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go-playground/cmd/service"
	"go-playground/configs"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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
		wantBody    map[string]string
	}{
		"no resources with get method and then 404": {
			inMethod:    "GET",
			inResources: "/",
			wantCode:    404,
			wantBody:    map[string]string{"detail": "/ resource does not exist", "instant": "/", "request_id": "", "title": "Resource Not Found", "type": "not:blank"}},
		"no resources with post method and then 404": {
			inMethod:    "POST",
			inResources: "/",
			wantCode:    404,
			wantBody:    map[string]string{"detail": "/ resource does not exist", "instant": "/", "request_id": "", "title": "Resource Not Found", "type": "not:blank"}},
		"no resources with put method and then 404": {
			inMethod:    "PUT",
			inResources: "/",
			wantCode:    404,
			wantBody:    map[string]string{"detail": "/ resource does not exist", "instant": "/", "request_id": "", "title": "Resource Not Found", "type": "not:blank"}},
		"no resources with delete method and then 404": {
			inMethod:    "DELETE",
			inResources: "/",
			wantCode:    404,
			wantBody:    map[string]string{"detail": "/ resource does not exist", "instant": "/", "request_id": "", "title": "Resource Not Found", "type": "not:blank"}},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			req, err := http.NewRequest(v.inMethod, testUriWith(v.inResources), nil)
			assert.NoError(t, err)
			res, err := client.Do(req)
			assert.NoError(t, err)

			var gotBody map[string]string
			err = json.NewDecoder(res.Body).Decode(&gotBody)
			assert.NoError(t, err)
			assert.Equal(t, v.wantCode, res.StatusCode)
			assert.Equal(t, v.wantBody, gotBody)
		})
	}
}

func TestStatusHandler(t *testing.T) {
	tests := map[string]struct {
		inMethod    string
		inResources string
		wantCode    int
		wantBody    map[string]string
	}{
		"/statuses and then success": {
			inMethod:    "GET",
			inResources: "/statuses",
			wantCode:    200,
			wantBody:    map[string]string{"status": "ok"}},
		"/statuses with post method and then 405": {
			inMethod:    "POST",
			inResources: "/statuses",
			wantCode:    405,
			wantBody:    map[string]string{"detail": "Http method POST is not allowed for /statuses resource", "instant": "/statuses", "request_id": "", "title": "Method Not Allowed", "type": "not:blank"}},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			req, err := http.NewRequest(v.inMethod, testUriWith(v.inResources), nil)
			assert.NoError(t, err)
			res, err := client.Do(req)
			assert.NoError(t, err)

			var gotBody map[string]string
			err = json.NewDecoder(res.Body).Decode(&gotBody)
			assert.NoError(t, err)
			assert.Equal(t, v.wantCode, res.StatusCode)
			assert.Equal(t, v.wantBody, gotBody)
		})
	}
}

func TestHelloHandler_GetName(t *testing.T) {
	tests := map[string]struct {
		inMethod    string
		inResources string
		inHeader    map[string]string
		wantCode    int
		wantBody    map[string]string
	}{
		"/hello with valid token and then success": {
			inMethod:    "GET",
			inResources: "/hello",
			inHeader:    map[string]string{"Authorization": "admin"}, wantCode: 200, wantBody: map[string]string{"message": "Hello tecchu11(ADMIN)!! You have Admin role."}},
		"/hello with invalid token and then 401": {
			inMethod:    "GET",
			inResources: "/hello",
			inHeader:    map[string]string{"Authorization": "no-auth"},
			wantCode:    401,
			wantBody:    map[string]string{"detail": "Request token was not found in your request header", "instant": "/hello", "request_id": "", "title": "Request With No Authentication", "type": "not:blank"}},
		"/hello with missing token and then 401": {
			inMethod:    "GET",
			inResources: "/hello",
			inHeader:    map[string]string{},
			wantCode:    401,
			wantBody:    map[string]string{"detail": "Request token was not found in your request header", "instant": "/hello", "request_id": "", "title": "Request With No Authentication", "type": "not:blank"}},
		"/hello with post method and then 405": {
			inMethod:    "POST",
			inResources: "/hello",
			inHeader:    map[string]string{"Authorization": "admin"},
			wantCode:    405,
			wantBody:    map[string]string{"detail": "Http method POST is not allowed for /hello resource", "instant": "/hello", "request_id": "", "title": "Method Not Allowed", "type": "not:blank"}},
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

			var gotBody map[string]string
			err = json.NewDecoder(res.Body).Decode(&gotBody)
			assert.NoError(t, err)
			assert.Equal(t, v.wantCode, res.StatusCode)
			assert.Equal(t, v.wantBody, gotBody)
		})
	}

}

func testUriWith(endpoint string) string {
	return svr.URL + endpoint
}
