package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ursuldaniel/brunoyam/internal/domain/models"
)

type Sqlite3Storage struct {
	db *sql.DB
}

func NewSqlite3Storage(database string) (*Sqlite3Storage, error) {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := CreateSqlite3DB(db); err != nil {
		return nil, err
	}

	return &Sqlite3Storage{
		db: db,
	}, nil
}

func CreateSqlite3DB(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	name TEXT,
    	email TEXT,
		password TEXT
	);
	
	CREATE TABLE IF NOT EXISTS tokens (
		token TEXT
	);
	`

	_, err := db.Exec(query)
	return err
}

func (s *Sqlite3Storage) ListUsers() ([]*models.User, error) {
	rows, err := s.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*models.User{}
	for rows.Next() {
		user := &models.User{}

		err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Password,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *Sqlite3Storage) GetUser(id int) (*models.User, error) {
	query := `SELECT * FROM users WHERE id = $1`

	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := &models.User{}
	for rows.Next() {
		err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Password,
		)

		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (s *Sqlite3Storage) CreateUser(user *models.User) error {
	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3)`

	_, err := s.db.Exec(query, user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sqlite3Storage) UpdateUser(id int, newUser *models.User) error {
	query := `UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4`

	_, err := s.db.Exec(query, newUser.Name, newUser.Email, newUser.Password, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sqlite3Storage) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sqlite3Storage) Login(loginUser *models.LoginUser) (int, error) {
	query := `SELECT id, password FROM users WHERE name = $1`

	rows, err := s.db.Query(query, loginUser.Name)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	var id int
	var password string
	for rows.Next() {
		err := rows.Scan(
			&id,
			&password,
		)

		if err != nil {
			return -1, err
		}
	}

	if password != loginUser.Password {
		return -1, fmt.Errorf("invalid data")
	}

	return id, nil
}
