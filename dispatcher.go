package main

import (
	"database/sql"
	"sync"
)

type Dispatcher struct {
	db        *sql.DB
	done      chan bool
	waitGroup *sync.WaitGroup
	tasks     map[int64]string
}

func NewDispatcher(db *sql.DB) *Dispatcher {
	d := &Dispatcher{
		db:        db,
		done:      make(chan bool),
		waitGroup: &sync.WaitGroup{},
	}
	d.waitGroup.Add(1)
	return d
}
