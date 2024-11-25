package storage

import "github.com/sayansaha934/students-api/internal/types"
type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
	DeleteStudentById(id int64) error
	UpdateStudentById(id int64, name string, email string, age int) error
}