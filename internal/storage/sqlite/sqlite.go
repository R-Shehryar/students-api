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
   