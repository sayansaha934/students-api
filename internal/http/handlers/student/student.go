package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sayansaha934/students-api/internal/storage"
	"github.com/sayansaha934/students-api/internal/types"
	"github.com/sayansaha934/students-api/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a student")
		var student types.Student
		err:=json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}
		if err!=nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		if err:= validator.New().Struct(student); err!=nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return 
		}
		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		slog.Info("User created successfully", slog.String("userId",fmt.Sprint(lastId)))
		if err!=nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}
		
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
		}

}


func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r) // Use Gorilla Mux to extract variables from the path
		idStr, ok := vars["id"]
		if !ok {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("missing id in path")))
			return
		}
		intId, err:=strconv.ParseInt(idStr, 10, 64)
		if err!=nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		slog.Info("Getting student", slog.String("id", fmt.Sprint(intId)))
		student, err:=storage.GetStudentById(intId)
		if err!=nil {
			slog.Error("error getting user", slog.String("id", fmt.Sprint(intId)))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, student)
		
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Getting all students")
		students, err := storage.GetStudents()
		if err!=nil {
			slog.Error("error getting users")
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, students)
	}
}

func Delete(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r) // Use Gorilla Mux to extract variables from the path
		idStr, ok := vars["id"]
		if !ok {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("missing id in path")))
			return
		}
		intId, err:=strconv.ParseInt(idStr, 10, 64)
		if err!=nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		slog.Info("Deleting student", slog.String("id", fmt.Sprint(intId)))
		err = storage.DeleteStudentById(intId)
		if err!=nil {
			slog.Error("error deleting user", slog.String("id", fmt.Sprint(intId)))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, "Student deleted successfully")
	}
}

func Update(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r) // Use Gorilla Mux to extract variables from the path
		idStr, ok := vars["id"]
		if !ok {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("missing id in path")))
			return
		}
		intId, err:=strconv.ParseInt(idStr, 10, 64)
		if err!=nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		slog.Info("Updating student", slog.String("id", fmt.Sprint(intId)))
		var student types.Student
		err = json.NewDecoder(r.Body).Decode(&student)
		if err!=nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		if err:= validator.New().Struct(student); err!=nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return 
		}
		err = storage.UpdateStudentById(intId, student.Name, student.Email, student.Age)
		if err!=nil {
			slog.Error("error updating user", slog.String("id", fmt.Sprint(intId)))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, "Student updated successfully")
	}
}