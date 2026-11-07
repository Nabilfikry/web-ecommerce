package shared

import (
	"encoding/json"
	"net/http"
)

type ApiMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type ApiResponse[T any] struct {
	Success bool     `json:"success"`
	Message string   `json:"message,omitempty"`
	Data    T        `json:"data,omitempty"`
	Meta    *ApiMeta `json:"meta,omitempty"`
	Error   any      `json:"error,omitempty"`
}

func JSON[T any](w http.ResponseWriter, status int, data T, message string, meta *ApiMeta) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := ApiResponse[T]{
		Success: status >= 200 && status < 300,
		Message: message,
		Data:    data,
		Meta:    meta,
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func Success[T any](w http.ResponseWriter, status int, data T, message string) {
	JSON(w, status, data, message, nil)
}

func SuccessWithMeta[T any](w http.ResponseWriter, data T, meta *ApiMeta) {
	JSON(w, http.StatusOK, data, "", meta)
}

type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

type ErrorPayload struct {
	Code    string        `json:"code,omitempty"`
	Details []ErrorDetail `json:"details,omitempty"`
}

func Error(w http.ResponseWriter, status int, message string, code string, details []ErrorDetail) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	var errContent any = message
	if code != "" || len(details) > 0 {
		errContent = ErrorPayload{
			Code:    code,
			Details: details,
		}
	}

	resp := ApiResponse[any]{
		Success: false,
		Error:   errContent,
	}

	_ = json.NewEncoder(w).Encode(resp)
}

