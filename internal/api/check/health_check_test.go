package check

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Response struct {
	Ok bool `json:"ok"`
}

func TestHealthCheck(t *testing.T) {
	srv := httptest.NewServer(HealthCheck(slog.Default()))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/health-check")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code should be 200")

	var response Response
	json.NewDecoder(resp.Body).Decode(&response)

	assert.Equal(t, true, response.Ok, "Response should be true")
}
