package main

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Option struct {
	ID    int    `json:"id"`
	Text  string `json:"text"`
	Votes int    `json:"votes"`
}

type Store struct {
	db *pgxpool.Pool
}

func NewStore(dsn string) (*Store, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return &Store{db: pool}, nil
}

func (s *Store) ListOptions(ctx context.Context) ([]Option, error) {
	rows, err := s.db.Query(ctx, `SELECT id,text,votes FROM options ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []Option
	for rows.Next() {
		var o Option
		if err := rows.Scan(&o.ID, &o.Text, &o.Votes); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil
}

func (s *Store) AddOption(ctx context.Context, text string) (Option, error) {
	var id int
	err := s.db.QueryRow(ctx,
		`INSERT INTO options(text) VALUES($1) RETURNING id`, text).Scan(&id)
	if err != nil {
		return Option{}, err
	}
	return Option{ID: id, Text: text, Votes: 0}, nil
}

func (s *Store) Vote(ctx context.Context, optionID int, fp string) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// проверяем лимит 2
	var count int
	if err = tx.QueryRow(ctx,
		`SELECT COUNT(*) FROM votes WHERE fingerprint=$1`, fp).Scan(&count); err != nil {
		return err
	}
	if count >= 2 {
		return errors.New("limit reached")
	}

	// добавляем голос
	_, err = tx.Exec(ctx,
		`INSERT INTO votes(option_id,fingerprint) VALUES($1,$2)`, optionID, fp)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `UPDATE options SET votes=votes+1 WHERE id=$1`, optionID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
