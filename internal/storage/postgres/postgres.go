package postgres

import (
	"context"

	pgx "github.com/jackc/pgx/v5"
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
	);
	
	CREATE TABLE IF NOT EXISTS books (
		b_id SERIAL PRIMARY KEY,
		lable TEXT,
		author TEXT,
		u_id INTEGER
	);`

	_, err := conn.Exec(ctx, query)
	return err
}
