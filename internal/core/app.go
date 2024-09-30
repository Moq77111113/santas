package core

type (
	App struct {
		Notifier *Notifier
	}


)

func NewApp() *App {
	return &App{
		Notifier: &Notifier{
			clients: make(map[chan string]struct{}),
			Add:     make(chan chan string),
			Remove:  make(chan chan string),
		},
	}
}