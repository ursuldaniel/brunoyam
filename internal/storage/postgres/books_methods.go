package postgres

import (
	"context"
	"time"

	"github.com/ursuldaniel/brunoyam/internal/domain/models"
)

func (s *PostgresStorage) ListBooks(id int) ([]*models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `SELECT * FROM books WHERE u_id = $1`
	rows, err := s.conn.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	books := []*models.Book{}
	for rows.Next() {
		book := &models.Book{}
		if err := rows.Scan(&book.BookId, &book.Lable, &book.Author, &book.UserId); err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}

func (s *PostgresStorage) AddBook(id int, book *models.Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Prepare(ctx, "insert book", "INSERT INTO books (lable, author, u_id) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "insert book", book.Lable, book.Author, id)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *PostgresStorage) AddBooks(id int, books []*models.Book) error {
	for _, book := range books {
		if err := s.AddBook(id, book); err != nil {
			return err
		}
	}

	return nil
}
