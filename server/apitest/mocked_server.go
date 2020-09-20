package apitest

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ratedemon/go-rest/api/helper"
	"github.com/stretchr/testify/require"
)

type MockedServer struct {
	R      *require.Assertions
	Router *mux.Router
	// Client *http.Client
}

func newRouter(routes []helper.Route) *mux.Router {
	r := mux.NewRouter()
	for _, route := range routes {
		r.HandleFunc(route.Path, helper.HandleWrapper(route.HandleFunc)).Methods(route.Method)
	}
	return r
}

func NewServer(t *testing.T, routes []helper.Route) *MockedServer {
	return &MockedServer{
		R:      require.New(t),
		Router: newRouter(routes),
	}
}

// Do makes a HTTP requests and returns response as JSON body
func (m *MockedServer) Do(method, path string, jsonBody string) (int, string, error) {
	// prepare HTTP request
	var reqBody io.Reader
	if len(jsonBody) != 0 {
		reqBody = bytes.NewBufferString(jsonBody)
	}
	req, err := http.NewRequest(method, path, reqBody)
	if err != nil {
		return 0, "", fmt.Errorf("failed to prepare request: %w", err)
	}

	w := httptest.NewRecorder()

	m.Router.ServeHTTP(w, req)

	return w.Code, w.Body.String(), nil
}

// DoFile makes a HTTP requests with files inside
func (m *MockedServer) DoFile(method, path string, body *bytes.Buffer, writer *multipart.Writer) (int, string, error) {
	req, err := http.NewRequest(method, path, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		return 0, "", fmt.Errorf("failed to prepare request: %w", err)
	}

	w := httptest.NewRecorder()

	m.Router.ServeHTTP(w, req)

	return w.Code, w.Body.String(), nil
}
