package main

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{clients: make(map[*websocket.Conn]bool)}
}

func (h *Hub) broadcast(msg any) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for c := range h.clients {
		if err := c.WriteJSON(msg); err != nil {
			c.Close()
			delete(h.clients, c)
		}
	}
}

func (h *Hub) handler(store *Store) http.HandlerFunc {
	up := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := up.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade error:", err)
			return
		}
		h.mu.Lock()
		h.clients[c] = true
		h.mu.Unlock()

		// отправляем текущие данные при подключении
		opts, err := store.ListOptions(r.Context())
		if err != nil {
			log.Println("list error:", err)
			return
		}
		if err := c.WriteJSON(opts); err != nil {
			log.Println("write error:", err)
		}
	}
}

func (h *Hub) notifyOptions(store *Store) {
	opts, err := store.ListOptions(context.Background())
	if err == nil {
		h.broadcast(opts)
	}
}
