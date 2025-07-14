package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
)

var validText = regexp.MustCompile(`^[A-Za-zА-Яа-яЁё ]{1,40}$`)

type Server struct {
	store *Store
	hub   *Hub
}

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.hub.handler(s.store))
	mux.HandleFunc("/options", s.getOptions)
	mux.HandleFunc("/vote", s.postVote)
	mux.HandleFunc("/myvotes", s.getMyVotes)
	mux.HandleFunc("/option", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			s.deleteOption(w, r)
		} else if r.Method == http.MethodPost {
			s.postOption(w, r)
		}
	})
	return mux
}

func (s *Server) getMyVotes(w http.ResponseWriter, r *http.Request) {
	fp := r.URL.Query().Get("fp")
	if fp == "" {
		http.Error(w, "missing fingerprint", 400)
		return
	}
	rows, err := s.store.db.Query(r.Context(),
		`SELECT option_id FROM votes WHERE fingerprint=$1`, fp)
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		rows.Scan(&id)
		ids = append(ids, id)
	}
	json.NewEncoder(w).Encode(ids)
}

func (s *Server) getOptions(w http.ResponseWriter, r *http.Request) {
	opts, err := s.store.ListOptions(r.Context())
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}
	json.NewEncoder(w).Encode(opts)
}

func (s *Server) postVote(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OptionID int    `json:"option_id"`
		Fp       string `json:"fp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad body", 400)
		return
	}
	if err := s.store.ToggleVote(r.Context(), req.OptionID, req.Fp); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	s.hub.notifyOptions(s.store)
	w.WriteHeader(204)
}

func (s *Server) postOption(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad body", 400)
		return
	}

	req.Text = strings.TrimSpace(req.Text)
	if !validText.MatchString(req.Text) {
		http.Error(w, "validation failed", 400)
		return
	}

	opt, err := s.store.AddOption(r.Context(), req.Text)
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}
	s.hub.notifyOptions(s.store)
	json.NewEncoder(w).Encode(opt)
}
