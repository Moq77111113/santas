package services

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/moq77111113/chmoly-santas/config"
)

type (
	SSEClient struct {
		Config   *config.Config
		mu       sync.RWMutex
		channels map[string][]chan string
	}

	Event struct {
		// ID of the event, set the EventSource last event ID value
		ID []byte
		// Data to send in the event
		Data []byte
		// Event type identifying the event. The browser can use this to decide how to handle the event
		Event []byte
		// Reconnection time, the browser will wait for the specified time before trying to reconnect (ms)
		Retry []byte
		// To send a comment to keep the connection alive
		Comment []byte
	}
)

func NewSSEClient(config *config.Config) *SSEClient {

	return &SSEClient{
		Config:   config,
		channels: make(map[string][]chan string),
	}
}

func (s *SSEClient) AddClient(c echo.Context, channel string) {

	mChan := make(chan string)

	s.mu.Lock()
	s.channels[channel] = append(s.channels[channel], mChan)
	s.mu.Unlock()

	w := c.Response()

	flusher, ok := w.Writer.(http.Flusher)
	if !ok {
		log.Println("ResponseWriter does not support Flusher")
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
			event := Event{
				Data: []byte(msg),
			}
			if err := event.SendTo(w); err != nil {
				return
			}
			flusher.Flush()
		case <-ticker.C:
			event := Event{
				Comment: []byte("keep-alive"),
			}
			if err := event.SendTo(w); err != nil {
				return
			}
			flusher.Flush()
		case <-c.Request().Context().Done():
			log.Printf("SSE client disconnected, ip: %s", c.RealIP())
			s.removeClient(channel, mChan)
			return
		}
	}
}

func (s *SSEClient) Broadcast(channelID, message string) {
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

func (s *SSEClient) removeClient(channelID string, mChan chan string) {
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

func (e *Event) SendTo(w io.Writer) error {
	if len(e.Data) == 0 && len(e.Comment) == 0 {
		return nil
	}

	if len(e.Data) > 0 {

		// Attempt to send the event ID
		if _, err := fmt.Fprintf(w, "id: %s\n", e.ID); err != nil {
			return err
		}

		// Split data by new line
		// Send each line as a separate data event
		sd := bytes.Split(e.Data, []byte("\n"))
		for _, d := range sd {
			if _, err := fmt.Fprintf(w, "data: %s\n", d); err != nil {
				return err
			}
		}

		// Send the event type
		if len(e.Event) > 0 {
			if _, err := fmt.Fprintf(w, "event: %s\n", e.Event); err != nil {
				return err
			}
		}

		// Send the reconnection time
		if len(e.Retry) > 0 {
			if _, err := fmt.Fprintf(w, "retry: %s\n", e.Retry); err != nil {
				return err
			}
		}
	}

	// Send a comment to keep the connection alive
	if len(e.Comment) > 0 {
		if _, err := fmt.Fprintf(w, ": %s\n", e.Comment); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprint(w, "\n"); err != nil {
		return err
	}

	return nil
}
