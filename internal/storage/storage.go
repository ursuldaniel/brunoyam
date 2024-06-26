package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ursuldaniel/brunoyam/internal/domain/models"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(database string) (*Storage, error) {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := CreateDB(db); err != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}

func CreateDB(db *sql.DB) error {
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

func (s *Storage) ListUsers() ([]*models.User, error) {
	rows, err := s.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*types.User{}
	for rows.Next() {
		user := &types.User{}

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

func (s *Storage) GetUser(id int) (*types.User, error) {
	query := `SELECT * FROM users WHERE id = $1`

	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := &types.User{}
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

func (s *Storage) CreateUser(user *types.User) error {
	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3)`

	_, err := s.db.Exec(query, user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateUser(id int, newUser *types.User) error {
	query := `UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4`

	_, err := s.db.Exec(query, newUser.Name, newUser.Email, newUser.Password, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Login(loginUser *types.User) (int, error) {
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
