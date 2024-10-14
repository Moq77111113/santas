package event

import (
	"bytes"
	"fmt"
	"io"
)

type (
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
