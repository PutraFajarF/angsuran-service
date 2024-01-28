package response

import (
	"encoding/json"
	"net/http"
)

type HttpResponse struct {
	Status     bool        `json:"status"`
	StatusCode string      `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

func HttpSuccessResponse(w http.ResponseWriter, status bool, code int, statuscode, message string, data interface{}) {
	resp := &HttpResponse{
		Status:     status,
		StatusCode: statuscode,
		Message:    message,
		Data:       data,
	}

	jsonResp, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonResp)
}

func HttpErrorResponse(w http.ResponseWriter, status bool, code int, statuscode, message string) {
	resp := &HttpResponse{
		Status:     status,
		StatusCode: statuscode,
		Message:    message,
	}

	jsonResp, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonResp)
}
