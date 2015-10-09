package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/andevery/autogram"
	"github.com/lib/pq"
	"log"
	"strings"
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
	tasks     map[int64]autogram.BackgroundTask
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
		tasks:     make(map[int64]autogram.BackgroundTask),
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

	d.Error(1, errors.New("asasasasasasasasas"))
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
	client, err := d.getClient(notify.Data.ProfilesID)
	if err != nil {
		log.Printf("Get client: %s\n", err)
		d.stopTask(notify)
	}
	switch notify.Data.Type {
	case UNFOLLOW:
		d.startUnfollowTask(notify, client)
	case TAGS:
		d.startTagsTask(notify, client)
	}

	d.clients.IncTasks(notify.Data.ProfilesID)

	_, err = d.db.Exec("UPDATE tasks SET status=$1 WHERE id=$2", PROGRESS, notify.Data.ID)
	if err != nil {
		log.Printf("DB Exec on start: %s\n", err)
	}
}

func (d *Dispatcher) startUnfollowTask(notify *Notify, client *autogram.Client) {
	whitelist, err := d.getWhitelist(notify)
	if err != nil {
		d.Fatal(notify.Data.ID, err)
	}
	worker := autogram.NewUnfollowWorker(notify.Data.ID, client, whitelist, d)
	worker.Start()
	d.tasks[notify.Data.ID] = worker
}

func (d *Dispatcher) getWhitelist(notify *Notify) ([]string, error) {
	var whitelist []string
	rows, err := d.db.Query("SELECT username FROM whitelists WHERE profilesid=$1", notify.Data.ProfilesID)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return []string{}, err
		}
		whitelist = append(whitelist, username)
	}
	if err := rows.Err(); err != nil {
		return []string{}, err
	}

	return whitelist, nil
}

func (d *Dispatcher) startTagsTask(notify *Notify, client *autogram.Client) {
	tags := strings.Split(notify.Data.Tags, ",")
	worker := autogram.DefaultTagsWorker(notify.Data.ID, client, tags, d)
	worker.Follow = notify.Data.Follows
	worker.Like = notify.Data.Likes
	worker.LikesPerUser.Min = notify.Data.MinLikes
	worker.LikesPerUser.Max = notify.Data.MaxLikes
	worker.Start()
	d.tasks[notify.Data.ID] = worker
}

func (d *Dispatcher) pauseTask(notify *Notify) {
	if task, ok := d.tasks[notify.Data.ID]; ok {
		task.Stop()
		delete(d.tasks, notify.Data.ID)
		d.clients.DecTasks(notify.Data.ProfilesID)
	}
	_, err := d.db.Exec("UPDATE tasks SET status=$1 WHERE id=$2", PAUSED, notify.Data.ID)
	if err != nil {
		log.Printf("DB Exec on pause: %s\n", err)
	}
}

func (d *Dispatcher) stopTask(notify *Notify) {
	if task, ok := d.tasks[notify.Data.ID]; ok {
		task.Stop()
		delete(d.tasks, notify.Data.ID)
		d.clients.DecTasks(notify.Data.ProfilesID)
	}
	_, err := d.db.Exec("UPDATE tasks SET status=$1 WHERE id=$2", FINISHED, notify.Data.ID)
	if err != nil {
		log.Printf("DB Exec on stop: %s\n", err)
	}
}

func (d *Dispatcher) getClient(id int64) (*autogram.Client, error) {
	client, ok := d.clients.Get(id)
	if ok {
		log.Println("Client received from pool")
		return client, nil
	}

	var username string
	var password string
	err := d.db.QueryRow("SELECT username, password FROM profiles WHERE id=$1", id).Scan(&username, &password)
	if err != nil {
		return nil, err
	}

	client, err = autogram.DefaultClient(username, password, "2079178474.1fb234f.682a311e35334df3842ccb654516baf5")
	if err != nil {
		return nil, err
	}

	d.clients.Add(id, client)
	log.Println("Client created and added to pool")
	return client, nil
}

func (d *Dispatcher) Report(id int64, report map[string]int) {
	followed, follow_ok := report["followed"]
	liked, like_ok := report["liked"]
	unfollowed, unfollow_ok := report["unfollowed"]
	if follow_ok && like_ok {
		_, err := d.db.Exec("UPDATE tasks SET followscount=$1, likescount=$2 WHERE id=$3", followed, liked, id)
		if err != nil {
			log.Printf("Report update: %s\n", err)
		}
	}
	if unfollow_ok {
		_, err := d.db.Exec("UPDATE tasks SET unfollowscount=$1 WHERE id=$2", unfollowed, id)
		if err != nil {
			log.Printf("Report update: %s\n", err)
		}
	}
}

func (d *Dispatcher) Error(id int64, err error) {
	d.db.Exec("INSERT INTO errors (type, message, tasksid) VALUES ($1, $2, $3)", "error", err.Error(), id)
}

func (d *Dispatcher) Fatal(id int64, err error) {
	d.db.Exec("INSERT INTO errors (type, message, tasksid) VALUES ($1, $2, $3)", "fatal", err.Error(), id)
	if task, ok := d.tasks[id]; ok {
		task.Stop()
		_, err := d.db.Exec("UPDATE tasks SET status=$1 WHERE id=$2", FINISHED, id)
		if err != nil {
			log.Printf("Fatal DB Exec on stop: %s\n", err)
		}
	}
}
