package backend

import (
	"encoding/json"
	"net/http"
)

func Ok200(v any, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	e := json.NewEncoder(w)
	e.Encode(map[string]any{
		"code":    200,
		"message": "Ok",
		"result":  v,
	})
}

func Error404Response(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	e := json.NewEncoder(w)
	e.Encode(map[string]any{
		"code":    404,
		"message": "Not Found",
		"result":  map[string]string{},
	})
}

func Error400Response(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	e := json.NewEncoder(w)
	e.Encode(map[string]any{
		"code":    400,
		"message": "Bad Request",
		"result":  map[string]string{},
	})
}

func Error500Response(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	e := json.NewEncoder(w)
	e.Encode(map[string]any{
		"code":    500,
		"message": "Internal Server Error",
		"result":  map[string]string{},
	})
}
