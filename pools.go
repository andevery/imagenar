package main

import (
	"github.com/andevery/autogram"
	"sync"
)

type Client struct {
	c          *autogram.Client
	tasksCount int
}

type ClientsPool struct {
	clients map[int64]*Client
	mutex   *sync.Mutex
}

func NewClientsPool() *ClientsPool {
	return &ClientsPool{
		clients: make(map[int64]*client),
		mutex:   &sync.Mutex,
	}
}

func (p *ClientsPool) Get(id int64) (*autogram.Client, ok) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	c, ok := p.clients[i]
	if !ok {
		return nil, false
	}

	return c.c, true
}

func (p *ClientsPool) Add(id int64, c *autogram.Client) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.clients[id] = &Client{c: c, tasksCount: 0}
}

func (p *ClientsPool) IncTasks(id int64) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.clients[id].tasksCount++
}

func (p *ClientsPool) DecTasks(id int64) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.clients[id].tasksCount--
	if p.clients[id].tasksCount == 0 {
		delete(p.clients, id)
	}
}
