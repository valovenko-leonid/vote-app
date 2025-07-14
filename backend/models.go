package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

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

func (s *Store) ToggleVote(ctx context.Context, optionID int, fp string) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// проверяем, голосует ли повторно (для снятия)
	var exists bool
	err = tx.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM votes WHERE option_id=$1 AND fingerprint=$2)`, optionID, fp).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		// снимаем голос
		_, err = tx.Exec(ctx,
			`DELETE FROM votes WHERE option_id=$1 AND fingerprint=$2`, optionID, fp)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `UPDATE options SET votes = votes - 1 WHERE id=$1`, optionID)
		if err != nil {
			return err
		}
	} else {
		// проверяем общее число голосов
		var count int
		err = tx.QueryRow(ctx,
			`SELECT COUNT(*) FROM votes WHERE fingerprint=$1`, fp).Scan(&count)
		if err != nil {
			return err
		}
		if count >= 2 {
			return errors.New("Вы уже выбрали 2 варианта")
		}

		// добавляем голос
		_, err = tx.Exec(ctx, `INSERT INTO votes(option_id, fingerprint) VALUES($1, $2)`, optionID, fp)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `UPDATE options SET votes = votes + 1 WHERE id=$1`, optionID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s *Server) deleteOption(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id", 400)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}

	// DELETE голосов
	_, err = s.store.db.Exec(r.Context(), `DELETE FROM votes WHERE option_id=$1`, id)
	if err != nil {
		log.Println("Ошибка удаления голосов:", err)
		http.Error(w, "db error (votes)", 500)
		return
	}

	// DELETE варианта
	_, err = s.store.db.Exec(r.Context(), `DELETE FROM options WHERE id=$1`, id)
	if err != nil {
		log.Println("Ошибка удаления варианта:", err)
		http.Error(w, "db error (option)", 500)
		return
	}

	s.hub.notifyOptions(s.store)
	w.WriteHeader(204)
}
