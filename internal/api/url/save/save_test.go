package save

import (
	"bytes"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) SaveURL(url string, alias string) (int64, error) {
	args := m.Called(url, alias)
	return args.Get(0).(int64), args.Error(1)
}

func TestSave(t *testing.T) {
	t.Run("Save with valid url", func(t *testing.T) {
		saver := new(MockStorage)
		saver.On("SaveURL", "https://example.com", "alias").Return(int64(1), nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(`{"url": "https://example.com", "alias": "alias"}`)))

		Save(saver, slog.Default()).ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("Wrong status code: %v", w.Code)
		}
	})

	t.Run("Save with invalid url", func(t *testing.T) {
		saver := new(MockStorage)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(`{"url": "invalid://example.com", "alias": "alias"}`)))

		Save(saver, slog.Default()).ServeHTTP(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Wrong status code: %v", w.Code)
		}
	})

	t.Run("Save with empty url", func(t *testing.T) {
		saver := new(MockStorage)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(`{"url": "", "alias": "alias"}`)))

		Save(saver, slog.Default()).ServeHTTP(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Wrong status code: %v", w.Code)
		}
	})
}
