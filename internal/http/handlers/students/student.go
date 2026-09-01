package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"log/slog"
	"strconv"

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
func  GetStudentByID(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		intId,err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GenreralErrorResponse(fmt.Errorf("invalid student ID")))
			return
		}
		slog.Info("Fetching student by ID.", slog.String("id", id))
		student, err := storage.GetStudentByID(intId)
		
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GenreralErrorResponse(err))
			return
		}
		response.WriteJSON(w, http.StatusOK, student)
	}

}

func GetAllStudents(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implement logic to fetch all students from storage
		slog.Info("Fetching all students.")
		students, err := storage.GetAllStudents()
		
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GenreralErrorResponse(err))
			return
		}
		response.WriteJSON(w, http.StatusOK, students)
	}
}

func  UpdateStudentByID(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student
		id := r.PathValue("id")
		intId,err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GenreralErrorResponse(fmt.Errorf("invalid student ID")))
			return
		}
		err = json.NewDecoder(r.Body).Decode(&student)
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
		updatedStudent, storageErr := storage.UpdateStudentByID(intId, student.Name, student.Email, student.Age)
		if storageErr != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GenreralErrorResponse(storageErr))
			return
		}
		response.WriteJSON(w, http.StatusOK, updatedStudent)

	}

}

func  DeleteStudentByID(storage storage.Storage) http.HandlerFunc {
   	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		intId,err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GenreralErrorResponse(fmt.Errorf("invalid student ID")))
			return
		}
		slog.Info("Fetching student by ID.", slog.String("id", id))
		err = storage.DeleteStudentByID(intId)
		
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GenreralErrorResponse(err))
			return
		}
		response.WriteJSON(w, http.StatusOK, map[string]interface{}{"message": "Student deleted successfully"})
	}

}