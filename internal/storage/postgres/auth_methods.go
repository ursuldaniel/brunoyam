package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/ursuldaniel/brunoyam/internal/domain/models"
)

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
