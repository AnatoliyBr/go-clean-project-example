package http

import (
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AnatoliyBr/go-clean-project-example/internal/repository/counterstore"
	"github.com/AnatoliyBr/go-clean-project-example/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestHTTPServerSet(t *testing.T) {
	testCases := []struct {
		name     string
		queryStr string
		want     string
	}{
		{
			name:     "normal",
			queryStr: "name=a&val=432",
			want:     "counter 'a' with val '432' was set\n",
		},
		{
			name:     "overflow",
			queryStr: "name=a&val=9223372036854775808",
			want:     "strconv.Atoi: parsing \"9223372036854775808\": value out of range\n",
		},
		{
			name:     "negative",
			queryStr: "name=a&val=-4",
			want:     "counter 'a' can't be negative\n",
		},
		{
			name:     "invalid",
			queryStr: "qwerty",
			want:     "strconv.Atoi: parsing \"\": invalid syntax\n",
		},
	}

	// Repository
	cs := counterstore.NewCounterStore(map[string]int{"i": 0, "j": 0})

	// Usecase
	uc := usecase.NewUseCase(cs)

	// Controller
	s := NewHTTPServer(NewConfig(), uc)

	go s.StartCounterManager()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// httptest.ResponseRecorder implements http.ResponseWriter interface
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/set?%s", tc.queryStr), nil)
			s.set(rec, req)
			assert.Equal(t, tc.want, rec.Body.String())
		})
	}
}

func TestHTTPServerGet(t *testing.T) {
	testCases := []struct {
		name     string
		queryStr string
		want     string
	}{
		{
			name:     "existent",
			queryStr: "name=i",
			want:     "'i': 0\n",
		},
		{
			name:     "nonexistent",
			queryStr: "name=a",
			want:     "'a' not found\n",
		},
	}

	// Repository
	cs := counterstore.NewCounterStore(map[string]int{"i": 0, "j": 0})

	// Usecase
	uc := usecase.NewUseCase(cs)

	// Controller
	s := NewHTTPServer(NewConfig(), uc)

	go s.StartCounterManager()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// httptest.ResponseRecorder implements http.ResponseWriter interface
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/get?%s", tc.queryStr), nil)
			s.get(rec, req)
			assert.Equal(t, tc.want, rec.Body.String())
		})
	}
}

func TestHTTPServerInc(t *testing.T) {
	testCases := []struct {
		name     string
		queryStr string
		want     string
	}{
		{
			name:     "normal",
			queryStr: "name=i",
			want:     "ok, 'i': 1\n",
		},
		{
			name:     "nonexistent",
			queryStr: "name=a",
			want:     "'a' not found\n",
		},
		{
			name:     "max",
			queryStr: "name=j",
			want:     "'j' has reached its maximum\n",
		},
	}

	// Repository
	cs := counterstore.NewCounterStore(map[string]int{"i": 0, "j": math.MaxInt})

	// Usecase
	uc := usecase.NewUseCase(cs)

	// Controller
	s := NewHTTPServer(NewConfig(), uc)

	go s.StartCounterManager()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// httptest.ResponseRecorder implements http.ResponseWriter interface
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/inc?%s", tc.queryStr), nil)
			s.inc(rec, req)
			assert.Equal(t, tc.want, rec.Body.String())
		})
	}
}

func TestHTTPServerDec(t *testing.T) {
	testCases := []struct {
		name     string
		queryStr string
		want     string
	}{
		{
			name:     "normal",
			queryStr: "name=i",
			want:     "ok, 'i': 9\n",
		},
		{
			name:     "nonexistent",
			queryStr: "name=a",
			want:     "'a' not found\n",
		},
		{
			name:     "negative",
			queryStr: "name=j",
			want:     "'j' can't be negative\n",
		},
	}

	// Repository
	cs := counterstore.NewCounterStore(map[string]int{"i": 10, "j": 0})

	// Usecase
	uc := usecase.NewUseCase(cs)

	// Controller
	s := NewHTTPServer(NewConfig(), uc)

	go s.StartCounterManager()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// httptest.ResponseRecorder implements http.ResponseWriter interface
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/dec?%s", tc.queryStr), nil)
			s.dec(rec, req)
			assert.Equal(t, tc.want, rec.Body.String())
		})
	}
}
