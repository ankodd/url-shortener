package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// Write wrapper for writing response
func Write(w http.ResponseWriter, status int, err error, data ...interface{}) {
	r := Response{
		Status: status,
	}

	if len(data) > 0 {
		r.Data = data[0]
	}

	if err != nil {
		r.Error = err.Error()
	}

	w.WriteHeader(r.Status)

	err = json.NewEncoder(w).Encode(&r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
