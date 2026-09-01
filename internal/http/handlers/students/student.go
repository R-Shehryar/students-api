package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"log/slog"

	"github.com/R-Shehryar/students-api/internal/types"
	"github.com/R-Shehryar/students-api/internal/utils/response"
	"github.com/R-Shehryar/students-api/internal/storage"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
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
		lastId, storageErr := storage.CreateStudent(student.Name, student.Email, student.Age)
		if storageErr != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GenreralErrorResponse(storageErr))
			return
		}
		student.ID = lastId
		slog.Info("Student created successfully.", slog.Int64("id", student.ID), slog.String("name", student.Name), slog.String("email", student.Email))
		response.WriteJSON(w, http.StatusCreated, student)
	}

}
