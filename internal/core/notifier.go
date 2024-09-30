package core

import "fmt"

type Notifier struct {
	clients map[chan string]struct{}
	Add     chan chan string
	Remove  chan chan string
}

func (n *Notifier) Notify(message string) {
	fmt.Printf("Notifying %v clients with: %v\n", len(n.clients), message)
	for client := range n.clients {
		client <- message
	}
}

func (n *Notifier) Listen() {
	for {
		select {
		case client := <-n.Add:
			n.clients[client] = struct{}{}
		case client := <-n.Remove:
			delete(n.clients, client)
		}
	}
}

func (n *Notifier) Close() {
	for client := range n.clients {
		close(client)
	}
}
