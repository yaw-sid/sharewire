package main

import (
	"context"
	"fmt"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(p Person) string {
	return fmt.Sprintf("Hello %s (Age: %d)!", p.Name, p.Age)
}

func (a *App) ListPeers() []Peer {
	return []Peer{
		*NewPeer("wed3e2e2w", "Peer 1", 1),
		*NewPeer("wed3e2e2w", "Peer 2", 2),
		*NewPeer("wed3e2e2w", "Peer 3", 3),
		*NewPeer("wed3e2e2w", "Peer 4", 4),
	}
}
