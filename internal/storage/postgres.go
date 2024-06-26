package storage

import (
	"context"
	"fmt"
	"time"

	pgx "github.com/jackc/pgx/v5"
	"github.com/ursuldaniel/brunoyam/internal/domain/models"
)

type PostgresStorage struct {
	conn *pgx.Conn
}

func NewPostgresStorage(ctx context.Context, connStr string) (*PostgresStorage, error) {
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}
	// defer conn.Close(context.Background())

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	if err := CreatePostgresDB(ctx, conn); err != nil {
		return nil, err
	}

	return &PostgresStorage{
		conn: conn,
	}, nil
}

func CreatePostgresDB(ctx context.Context, conn *pgx.Conn) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		password TEXT	
	);`

	_, err := conn.Exec(ctx, query)
	return err
}

func (s *PostgresStorage) ListUsers() ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	rows, err := s.conn.Query(ctx, "SELECT * FROM users")
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

func (s *PostgresStorage) GetUser(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `SELECT * FROM users WHERE id = $1`
	rows, err := s.conn.Query(ctx, query, id)
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

func (s *PostgresStorage) CreateUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3)`
	_, err := s.conn.Exec(ctx, query, user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) UpdateUser(id int, newUser *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4`
	_, err := s.conn.Exec(ctx, query, newUser.Name, newUser.Email, newUser.Password, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) DeleteUser(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `DELETE FROM users WHERE id = $1`
	_, err := s.conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) Login(loginUser *models.LoginUser) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `SELECT id, password FROM users WHERE name = $1`
	rows, err := s.conn.Query(ctx, query, loginUser.Name)
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
