package sqlite

import (
	"database/sql"

	"github.com/R-Shehryar/students-api/internal/config"
	"github.com/R-Shehryar/students-api/internal/types"
	_ "modernc.org/sqlite"

	"fmt"
)

type SQLite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*SQLite, error) {
	db, err := sql.Open("sqlite", cfg.StoragePath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		age INTEGER
	)`)
	if err != nil {
		return nil, err
	}
	return &SQLite{Db: db}, nil
}
func (s *SQLite) CreateStudent(name string, email string, age int) (int64, error) {
    stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
    if err != nil {
        return 0, err
    }
    defer stmt.Close()
    result, err := stmt.Exec(name, email, age)
    if err != nil {
        return 0, err
    }
    lastId, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }
    return lastId, nil
}

func (s *SQLite) GetStudentByID(id int64) (types.Student, error) {
	
	stmt, err := s.Db.Prepare("SELECT * FROM students WHERE id = ?")
	if err != nil {
        return types.Student{}, err
    }
    defer stmt.Close()
	 var student types.Student
     err = stmt.QueryRow(id).Scan(&student.ID, &student.Name, &student.Email, &student.Age)
    
    if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("student with ID %d not found", id)
		}
        return types.Student{}, fmt.Errorf("failed to get student by ID: %v", err)
    }
    return student, nil
}
   
func (s *SQLite) GetAllStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM students")
	if err != nil {
        return []types.Student{}, err
    }
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return []types.Student{}, err
	}
	defer rows.Close()
	var students []types.Student
	for rows.Next() {
		var student types.Student
		err = rows.Scan(&student.ID, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return []types.Student{}, err
		}
		students = append(students, student)
	}
	return students, nil
}
func (s *SQLite) UpdateStudentByID(id int64, name string, email string, age int) (types.Student, error) {
	stmt, err := s.Db.Prepare("UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age, id)
	if err != nil {
		return types.Student{}, fmt.Errorf("failed to update student: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return types.Student{}, fmt.Errorf("failed to check update: %v", err)
	}

	if rowsAffected == 0 {
		return types.Student{}, fmt.Errorf("student with ID %d not found", id)
	}

	return types.Student{
		ID:    id,
		Name:  name,
		Email: email,
		Age:   age,
	}, nil
}
func (s *SQLite) DeleteStudentByID(id int64) error {
	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to delete student: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("student with ID %d not found", id)
	}

	return nil
}