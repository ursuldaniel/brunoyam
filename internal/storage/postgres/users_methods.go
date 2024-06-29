package postgres

import (
	"context"
	"time"

	"github.com/ursuldaniel/brunoyam/internal/domain/models"
)

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

	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Prepare(ctx, "insert user", "INSERT INTO users (name, email, password) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "insert user", user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *PostgresStorage) AddUsers(users []*models.User) error {
	for _, user := range users {
		if err := s.CreateUser(user); err != nil {
			return err
		}
	}

	return nil
}

func (s *PostgresStorage) UpdateUser(id int, newUser *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Prepare(ctx, "update", "UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4")
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "update", newUser.Name, newUser.Email, newUser.Password, id)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
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
