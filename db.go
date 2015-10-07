package main

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"time"
)

type dbConf map[string]struct {
	Driver string `yaml:"driver"`
	Conn   string `yaml:"open"`
}

var (
	db            *sql.DB
	tasksListener *pq.Listener
)

func dbConnect(conf dbConf) (err error) {
	c, ok := conf[env()]
	if !ok {
		return errors.New("Configuration for \"" + env() + "\" environment not found.")
	}

	db, err = sql.Open(c.Driver, c.Conn)
	if err != nil {
		return
	}

	tasksListener = pq.NewListener(c.Conn, 10*time.Second, time.Minute, nil)
	err = tasksListener.Listen("tasks_notify_event")

	return
}
