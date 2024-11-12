//go:build go1.22

// Package oapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package oapi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Error defines model for Error.
type Error struct {
	// Message error message
	Message string `json:"message"`
}

// Simple defines model for Simple.
type Simple struct {
	// Message message
	Message string `json:"message"`
}

// Task defines model for Task.
type Task struct {
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	ID        string    `json:"id"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TaskContent defines model for TaskContent.
type TaskContent struct {
	// Content Content of task. Content must be not blank.
	Content string `json:"content"`
}

// Limit pagination limit size.
type Limit = int32

// Next pagination cursor value.
type Next = string

// TaskID ID of task.
type TaskID = string

// Response400 defines model for Response400.
type Response400 = Error

// Response404 defines model for Response404.
type Response404 = Error

// Response500 defines model for Response500.
type Response500 = Error

// ResponseHealthCheck defines model for ResponseHealthCheck.
type ResponseHealthCheck = Simple

// ResponseTask defines model for ResponseTask.
type ResponseTask = Task

// ResponseTaskID defines model for ResponseTaskID.
type ResponseTaskID struct {
	// ID ID of task.
	ID string `json:"id"`
}

// ResponseTasks defines model for ResponseTasks.
type ResponseTasks struct {
	// HasNext whether has next items.
	HasNext bool `json:"hasNext"`

	// Items Items of task
	Items []Task `json:"items"`

	// Next cursor of next item.
	Next string `json:"next"`
}

// ResponseUserID defines model for ResponseUserID.
type ResponseUserID struct {
	// ID ID of user id
	ID openapi_types.UUID `json:"id"`
}

// RequestTask defines model for RequestTask.
type RequestTask = TaskContent

// RequestUser defines model for RequestUser.
type RequestUser struct {
	// Email user email
	Email openapi_types.Email `json:"email"`

	// EmailVerified whether email is verified
	EmailVerified bool `json:"emailVerified"`

	// FamilyName user family name
	FamilyName string `json:"familyName"`

	// GivenName user given name
	GivenName string `json:"givenName"`
}

// ListTasksParams defines parameters for ListTasks.
type ListTasksParams struct {
	Next  *Next  `form:"next,omitempty" json:"next,omitempty"`
	Limit *Limit `form:"limit,omitempty" json:"limit,omitempty"`
}

// PostUserJSONBody defines parameters for PostUser.
type PostUserJSONBody struct {
	// Email user email
	Email openapi_types.Email `json:"email"`

	// EmailVerified whether email is verified
	EmailVerified bool `json:"emailVerified"`

	// FamilyName user family name
	FamilyName string `json:"familyName"`

	// GivenName user given name
	GivenName string `json:"givenName"`
}

// PostTaskJSONRequestBody defines body for PostTask for application/json ContentType.
type PostTaskJSONRequestBody = TaskContent

// PutTaskJSONRequestBody defines body for PutTask for application/json ContentType.
type PutTaskJSONRequestBody = TaskContent

// PostUserJSONRequestBody defines body for PostUser for application/json ContentType.
type PostUserJSONRequestBody PostUserJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Health check API
	// (GET /health)
	HealthCheck(w http.ResponseWriter, r *http.Request)
	// List tasks
	// (GET /tasks)
	ListTasks(w http.ResponseWriter, r *http.Request, params ListTasksParams)
	// Post task
	// (POST /tasks)
	PostTask(w http.ResponseWriter, r *http.Request)
	// Get task
	// (GET /tasks/{taskId})
	GetTask(w http.ResponseWriter, r *http.Request, taskID TaskID)
	// Put task
	// (PUT /tasks/{taskId})
	PutTask(w http.ResponseWriter, r *http.Request, taskID TaskID)
	// Post user
	// (POST /users)
	PostUser(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// HealthCheck operation middleware
func (siw *ServerInterfaceWrapper) HealthCheck(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.HealthCheck(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// ListTasks operation middleware
func (siw *ServerInterfaceWrapper) ListTasks(w http.ResponseWriter, r *http.Request) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params ListTasksParams

	// ------------- Optional query parameter "next" -------------

	err = runtime.BindQueryParameter("form", true, false, "next", r.URL.Query(), &params.Next)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "next", Err: err})
		return
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "limit", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ListTasks(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PostTask operation middleware
func (siw *ServerInterfaceWrapper) PostTask(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostTask(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetTask operation middleware
func (siw *ServerInterfaceWrapper) GetTask(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "taskId" -------------
	var taskID TaskID

	err = runtime.BindStyledParameterWithOptions("simple", "taskId", r.PathValue("taskId"), &taskID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "taskId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetTask(w, r, taskID)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PutTask operation middleware
func (siw *ServerInterfaceWrapper) PutTask(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "taskId" -------------
	var taskID TaskID

	err = runtime.BindStyledParameterWithOptions("simple", "taskId", r.PathValue("taskId"), &taskID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "taskId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PutTask(w, r, taskID)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PostUser operation middleware
func (siw *ServerInterfaceWrapper) PostUser(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostUser(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{})
}

// ServeMux is an abstraction of http.ServeMux.
type ServeMux interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type StdHTTPServerOptions struct {
	BaseURL          string
	BaseRouter       ServeMux
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, m ServeMux) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseRouter: m,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, m ServeMux, baseURL string) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseURL:    baseURL,
		BaseRouter: m,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options StdHTTPServerOptions) http.Handler {
	m := options.BaseRouter

	if m == nil {
		m = http.NewServeMux()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	m.HandleFunc("GET "+options.BaseURL+"/health", wrapper.HealthCheck)
	m.HandleFunc("GET "+options.BaseURL+"/tasks", wrapper.ListTasks)
	m.HandleFunc("POST "+options.BaseURL+"/tasks", wrapper.PostTask)
	m.HandleFunc("GET "+options.BaseURL+"/tasks/{taskId}", wrapper.GetTask)
	m.HandleFunc("PUT "+options.BaseURL+"/tasks/{taskId}", wrapper.PutTask)
	m.HandleFunc("POST "+options.BaseURL+"/users", wrapper.PostUser)

	return m
}
