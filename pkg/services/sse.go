package services

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/moq77111113/chmoly-santas/config"
	"github.com/moq77111113/chmoly-santas/ent"
	"github.com/moq77111113/chmoly-santas/pkg/event"
)

type (
	SSEClient struct {
		Config   *config.Config
		mu       sync.RWMutex
		channels map[string][]chan Message
	}

	Message struct {
		Type string `json:"type"`
		Data string `json:"data"`
	}
)

func NewSSEClient(config *config.Config) *SSEClient {

	return &SSEClient{
		Config:   config,
		channels: make(map[string][]chan Message),
	}
}

func (s *SSEClient) AddClient(c echo.Context, channel string) {

	me, ok := c.Get("me").(*ent.Member)
	if !ok || me == nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}
	mChan := make(chan Message)

	s.mu.Lock()
	s.channels[channel] = append(s.channels[channel], mChan)
	s.mu.Unlock()

	w := c.Response()

	flusher, ok := w.Writer.(http.Flusher)
	if !ok {
		c.Logger().Error("ResponseWriter does not support Flusher")
		c.String(http.StatusInternalServerError, "Streaming unsupported!")
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case msg := <-mChan:

			bytes, err := msg.toBytesString()

			if err != nil {
				c.Logger().Errorf("Error marshalling message: %v", err)
				continue
			}
			event := event.Event{
				Data: []byte(bytes),
			}

			if err := event.SendTo(w); err != nil {
				c.Logger().Errorf("Error sending event: %v", err)
				return
			}
			flusher.Flush()
		case <-ticker.C:

			bytes, err := Message{
				Type: "keep-alive",
				Data: "keep-alive",
			}.toBytesString()

			if err != nil {
				c.Logger().Errorf("Error marshalling message: %v", err)
				continue
			}
			event := event.Event{
				Comment: []byte(bytes),
			}
			if err := event.SendTo(w); err != nil {
				return
			}
			flusher.Flush()
		case <-c.Request().Context().Done():
			log.Printf("SSE client disconnected, %s@%s", me.Name, c.RealIP())
			s.removeClient(channel, mChan)
			return
		}
	}
}

func (s *SSEClient) Broadcast(channelID string, message Message) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	clients, ok := s.channels[channelID]
	if !ok {
		return
	}

	for _, ch := range clients {
		ch <- message
	}
}

func (s *SSEClient) removeClient(channelID string, mChan chan Message) {
	s.mu.Lock()
	defer s.mu.Unlock()

	clients := s.channels[channelID]
	for i, ch := range clients {
		if ch == mChan {
			s.channels[channelID] = append(clients[:i], clients[i+1:]...)
			break
		}
	}

	if len(s.channels[channelID]) == 0 {
		delete(s.channels, channelID)
	}

	close(mChan)
}

func (m Message) toBytesString() ([]byte, error) {

	return json.Marshal(m)
}
