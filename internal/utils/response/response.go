package response

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func GenreralErrorResponse(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}
func ValidationErrorResponse(err validator.ValidationErrors) Response {
	var errMsg []string
	for _, fieldErr := range err {
		switch fieldErr.ActualTag() {
		case "required":
			errMsg = append(errMsg, fieldErr.Field()+" is required")
		default:
			errMsg = append(errMsg, fieldErr.Error()+" is invalid")
		}
	}
	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsg, ", "),
	}
}
