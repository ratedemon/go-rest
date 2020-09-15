package apitest

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ratedemon/go-rest/api/helper"
	"github.com/stretchr/testify/require"
)

type MockedServer struct {
	R      *require.Assertions
	Server *httptest.Server
	Client *http.Client
}

// Close closes underlying HTTP server
func (m *MockedServer) Close() {
	m.Server.Close()
}

// NewHandler returns a new testing handler for the provided HTTP routes
func NewHandler(t *testing.T, routes []helper.Route) http.Handler {
	r := http.NewServeMux()
	for _, route := range routes {
		r.HandleFunc(route.Path, helper.HandleWrapper(route.HandleFunc))
	}
	return r
}

func NewServer(t *testing.T, routes []helper.Route) *httptest.Server {
	return httptest.NewServer(NewHandler(t, routes))
}

// NewMockedServer creates new gateway helper
func NewMockedServer(t *testing.T, server *httptest.Server) *MockedServer {
	return &MockedServer{
		R:      require.New(t),
		Server: server,
		Client: http.DefaultClient,
	}
}

// NewMockedServerRoutes creates new gateway helper with routes
func NewMockedServerRoutes(t *testing.T, routes []helper.Route) *MockedServer {
	return NewMockedServer(t, NewServer(t, routes))
}

// Do makes a HTTP requests and returns response as JSON body
func (m *MockedServer) Do(method, path string, jsonBody string) (int, string, error) {
	// prepare HTTP request
	var reqBody io.Reader
	if len(jsonBody) != 0 {
		reqBody = bytes.NewBufferString(jsonBody)
	}
	req, err := http.NewRequest(method, m.Server.URL+path, reqBody)
	if err != nil {
		return 0, "", fmt.Errorf("failed to prepare request: %w", err)
	}

	// do HTTP request
	resp, err := m.Client.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	// read HTTP response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, "", fmt.Errorf("failed to read response: %w", err)
	}

	return resp.StatusCode, string(data), nil
}
