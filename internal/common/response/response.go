package response

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fanzru/bythen/internal/app/user/port/genhttp"
)

// BaseResponse struct
type BaseResponse struct {
	Code       string      `json:"code"`
	Message    string      `json:"message,omitempty"`
	TraceError string      `json:"traceError,omitempty" `
	Data       interface{} `json:"data,omitempty"`
	ServerTime int64       `json:"serverTime"`
	RequestID  string      `json:"requestId"`
	Errors     interface{} `json:"errors,omitempty"`
}

// Helper method to write a success response
func WriteSuccessResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	response := BaseResponse{
		Code:       "SUCCESS",
		Data:       data,
		ServerTime: time.Now().Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// Helper method to write an error response
func WriteErrorResponse(w http.ResponseWriter, message string, err error, statusCode int) {
	response := &genhttp.Error{
		Code:       "ERROR",
		Message:    err.Error(),
		ServerTime: int(time.Now().Unix()),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
