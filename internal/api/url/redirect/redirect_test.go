package redirect

import (
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockStorage struct {
	mock.Mock
}

func (s *MockStorage) GetUrl(alias string) (string, error) {
	args := s.Called(alias)
	return args.String(0), args.Error(1)
}

func TestGetUrl(t *testing.T) {
	t.Run("GetUrl with valid alias", func(t *testing.T) {
		getter := new(MockStorage)
		getter.On("GetUrl", "alias").Return("https://example.com", nil)

		url, err := getter.GetUrl("alias")
		if err != nil {
			t.Fatalf("GetUrl returned error: %v", err)
		}

		if url != "https://example.com" {
			t.Fatalf("GetUrl returned wrong url: %v", url)
		}
	})

	t.Run("GetUrl with invalid alias", func(t *testing.T) {
		getter := new(MockStorage)
		getter.On("GetUrl", "non-valid-alias").Return("", errors.New(AliasNotFound))

		url, err := getter.GetUrl("non-valid-alias")
		if err == nil {
			t.Fatalf("GetUrl not returned error: %v", err)
		}

		if url != "" {
			t.Fatalf("GetUrl returned wrong url: %v", url)
		}
	})
}

func TestRedirect(t *testing.T) {
	t.Run("Redirect with valid alias", func(t *testing.T) {
		getter := new(MockStorage)
		getter.On("GetUrl", "alias").Return("https://example.com", nil)

		r := mux.NewRouter()
		r.Handle("/test_redirect/{alias}", Redirect(getter, slog.Default())).Methods(http.MethodGet)

		srv := httptest.NewServer(r)
		defer srv.Close()

		resp, err := http.Get(srv.URL + "/test_redirect/alias")
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Wrong status code: %v", resp.StatusCode)
		}
	})

	t.Run("Redirect with invalid alias", func(t *testing.T) {
		getter := new(MockStorage)
		getter.On("GetUrl", "non-valid-alias").Return("", errors.New(AliasNotFound))

		r := mux.NewRouter()
		r.Handle("/test_redirect/{alias}", Redirect(getter, slog.Default())).Methods(http.MethodGet)

		srv := httptest.NewServer(r)
		defer srv.Close()

		resp, err := http.Get(srv.URL + "/non-valid-alias")
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("Wrong status code: %v", resp.StatusCode)
		}
	})
}
