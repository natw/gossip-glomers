package main

import (
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type App struct {
	msgLock      sync.Mutex
	neighborLock sync.Mutex
	logLock      sync.Mutex
	Node         *maelstrom.Node
	messages     map[int]bool
	neighbors    []string
}

func NewApp(node *maelstrom.Node) *App {
	a := App{
		Node:      node,
		messages:  make(map[int]bool),
		neighbors: []string{},
	}
	return &a
}

func (a *App) Logf(format string, v ...any) {
	a.logLock.Lock()
	defer a.logLock.Unlock()
	log.Printf(format, v...)
}
