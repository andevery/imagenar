package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/lib/pq"
	"log"
	"sync"
	"time"
)

const (
	STARTING = iota
	WAITING
	PROGRESS
	PAUSE
	PAUSED
	STOPPING
	FINISHED
)

const (
	UNFOLLOW = iota
	TAGS
)

type DBConf map[string]struct {
	Driver string `yaml:"driver"`
	Conn   string `yaml:"open"`
}

type Notify struct {
	Action string `json:"action"`
	Table  string `json:"table"`
	Data   struct {
		ID         int64  `json:"id"`
		Type       int    `json:"type"`
		Status     int    `json:"status"`
		ProfilesID int64  `json:"profilesid"`
		Tags       string `json:"tags"`

		MaxTags int `json:"maxtags"`

		Follows       bool `json:"follows"`
		MaxFollowedBy int  `json:"maxfollowedby"`
		MinFollowedBy int  `json:"minfollowedby"`
		MaxFollows    int  `json:"maxfollows"`
		MinFollows    int  `json:"minfollows"`
		MinMedia      int  `json:"minmedia"`

		Likes    bool `json:"likes"`
		MaxLikes int  `json:"maxlikes"`
		MinLikes int  `json:"minlikes"`

		FollowsCount   int `json:"followscount"`
		LikesCount     int `json:"likescount"`
		UnfollowsCount int `json:"unfollowscount"`
		Delay          int `json:"delay"`
	} `json:"data"`
}

type Dispatcher struct {
	db        *sql.DB
	listener  *pq.Listener
	waitGroup *sync.WaitGroup
	done      chan bool
	tasks     map[int64]string
	clients   *ClientsPool
}

func NewDispatcher(conf DBConf) (d *Dispatcher, err error) {
	c, ok := conf[env()]
	if !ok {
		err = errors.New("Configuration for \"" + env() + "\" environment not found.")
		return
	}

	d = &Dispatcher{
		done:      make(chan bool),
		waitGroup: &sync.WaitGroup{},
		listener:  pq.NewListener(c.Conn, 10*time.Second, time.Minute, nil),
		clients:   NewClientsPool(),
	}

	d.db, err = sql.Open(c.Driver, c.Conn)
	if err != nil {
		return
	}
	return
}

func (d *Dispatcher) Start() (err error) {
	err = d.listener.Listen("tasks_notify_event")
	if err != nil {
		return
	}

	go d.perform()
	// d.startExistsTasks()

	return
}

func (d *Dispatcher) Stop() {
	close(d.done)
	d.waitGroup.Wait()
}

func (d *Dispatcher) perform() {
	for {
		select {
		case notify := <-d.listener.Notify:
			go d.processNotify(notify)
		case <-d.done:
			return
		}
	}
}

func (d *Dispatcher) processNotify(notify *pq.Notification) {
	log.Println("Received data from channel [", notify.Channel, "] :")
	log.Println(notify.Extra)
	var n Notify
	err := json.Unmarshal([]byte(notify.Extra), &n)
	if err != nil {
		log.Printf("Unmarshal notify: %s\n%s", err, notify.Extra)
	}
	if n.Table == "tasks" && n.Action != "DELETE" {
		d.processTask(&n)
	}
}

func (d *Dispatcher) processTask(notify *Notify) {
	switch notify.Data.Status {
	case STARTING:
		d.startTask(notify)
	case PAUSE:
		d.pauseTask(notify)
	case STOPPING:
		d.stopTask(notify)
	}
}

func (d *Dispatcher) startTask(notify *Notify) {
	client := d.getClient(notify.Data.ProfilesID)
	switch notify.Data.Type {
	case UNFOLLOW:
		log.Println("Started unfollow task")
	case TAGS:
		log.Println("Started follow/liking task")
	}

	_, err := d.db.Exec("UPDATE tasks SET status=$1 WHERE id=$2", PROGRESS, notify.Data.ID)
	if err != nil {
		log.Printf("DB Exec on start: %s\n", err)
	}
}

func (d *Dispatcher) pauseTask(notify *Notify) {
	_, err := d.db.Exec("UPDATE tasks SET status=$1 WHERE id=$2", PAUSED, notify.Data.ID)
	if err != nil {
		log.Printf("DB Exec on pause: %s\n", err)
	}
}

func (d *Dispatcher) stopTask(notify *Notify) {
	_, err := d.db.Exec("UPDATE tasks SET status=$1 WHERE id=$2", FINISHED, notify.Data.ID)
	if err != nil {
		log.Printf("DB Exec on stop: %s\n", err)
	}
}

func (d *Dispatcher) getClient(id int64) *autogram.Client {

}
