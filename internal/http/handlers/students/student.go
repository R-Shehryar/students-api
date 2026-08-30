package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/R-Shehryar/students-api/internal/types"
	"github.com/R-Shehryar/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GenreralErrorResponse(fmt.Errorf("empty body")))
			return
		}
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GenreralErrorResponse(err))
			return
		}
		validationErr := validator.New().Struct(student)

		if validationErr != nil {
			validationErrs := validationErr.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationErrorResponse(validationErrs))
			return
		}
		response.WriteJSON(w, http.StatusCreated, student)
	}

}
