package main

import (
	"encoding/json"
	"fmt"
	"github.com/andevery/instax"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

var (
	count = 0
)

type Liker struct {
	Min            int
	Max            int
	MinBreak       time.Duration
	MaxBreak       time.Duration
	RateLimitPause time.Duration
	FeedDepth      int
	UserCond       struct {
		FollowedBy int
		Follows    int
		Media      int
	}
	MediaCond struct {
		TagsCount int
	}

	tags      []string
	client    *instax.Client
	webCLient *WebCLient

	breakTime   time.Duration
	counter     int
	depth       int
	likesNumber int
}

func NewLiker(tags []string, client *instax.Client, webClient *WebCLient) *Liker {
	liker := new(Liker)
	liker.client = client
	liker.webCLient = webClient
	liker.tags = tags

	liker.Min = 50
	liker.Max = 100
	liker.MinBreak = 30 * time.Minute
	liker.MaxBreak = 50 * time.Minute
	liker.RateLimitPause = 20 * time.Minute
	liker.FeedDepth = 10
	liker.UserCond.FollowedBy = 500
	liker.UserCond.Follows = 200
	liker.UserCond.Media = 50
	liker.MediaCond.TagsCount = 10

	return liker
}

func (l *Liker) flushCounts() {
	rand.Seed(time.Now().Unix())
	l.likesNumber = rand.Intn(l.Max-l.Min) + l.Min
	l.breakTime = time.Duration(rand.Int63n(int64(l.MaxBreak)-int64(l.MinBreak)) + int64(l.MinBreak))
	l.counter = 0
}

func (l *Liker) isUserMatch(user *instax.User) bool {
	return user.Counts.FollowedBy <= l.UserCond.FollowedBy &&
		user.Counts.Follows <= l.UserCond.Follows &&
		user.Counts.Media >= l.UserCond.Media
}

func (l *Liker) fakeLoadImages(media []instax.Media) {
	for i, _ := range media {
		http.Get(media[i].Images.Thumbnail.URL)
	}
}

func (l *Liker) checkAndLike(media []instax.Media) {
	var user *instax.User
	var err error

	if l.client.Limit() < 500 {
		fmt.Println("")
		log.Printf("Rate limit pause. Resume after %s.", l.RateLimitPause)
		time.Sleep(l.RateLimitPause)
	}

	for i, _ := range media {
		if media[i].UserHasLiked || len(media[i].Tags) > l.MediaCond.TagsCount {
			continue
		}

		user, err = l.client.User(media[i].User.ID)
		if err != nil {
			panic(err)
			return
		}

		if l.isUserMatch(user) {
			http.Get(media[i].Images.LowResolution.URL)
			l.webCLient.Like(&media[i])
			fmt.Print(" ✔ ")
			if rand.Intn(3) <= 1 {
				l.webCLient.Follow(user)
				fmt.Print(" + ")
			}
			l.counter++
			if l.counter >= l.likesNumber {
				fmt.Println("")
				log.Printf("Break time. Liked %v photos. Resume after %s", l.counter, l.breakTime)
				time.Sleep(l.breakTime)
				l.flushCounts()
			}
		}
	}
}

func (l *Liker) Start() {
	l.flushCounts()
	for _, tag := range l.tags {
		feed := l.client.MediaByTag(tag)
		media, err := feed.Get()
		if err != nil {
			panic(err)
		}
		l.depth = 1
		l.checkAndLike(media)
		for {
			l.depth++
			if !feed.CanNext() || l.depth > l.FeedDepth {
				fmt.Println("")
				log.Printf("End of feed: #%s.", tag)
				break
			}
			media, err = feed.Next()
			if err != nil {
				panic(err)
			}
			l.checkAndLike(media)
		}
	}
}

type DelayerFunc func() time.Duration

type WebCLient struct {
	Delayer   DelayerFunc
	UserAgent string
	Cookie    string
	CSRFToken string
	client    *http.Client
}

type Response struct {
	Status string `json:"status"`
}

func NewWebClient() *WebCLient {
	return &WebCLient{
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.85 Safari/537.36",
		client:    new(http.Client),
	}
}

func (c *WebCLient) do(method, path, referer string) error {
	if c.Delayer != nil {
		d := c.Delayer()
		time.Sleep(d)
	}

	u := &url.URL{
		Scheme: "https",
		Host:   "instagram.com",
		Path:   fmt.Sprintf("/web/%s", path),
	}

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Cookie", c.Cookie)
	req.Header.Add("Referer", referer)
	req.Header.Add("User-Agent", c.UserAgent)
	req.Header.Add("X-CSRFToken", c.CSRFToken)
	req.Header.Add("X-Instagram-AJAX", "1")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	resp, err := c.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var status Response
	err = json.Unmarshal(body, &status)
	if err != nil {
		fmt.Println("")
		log.Println(string(body))
		time.Sleep(time.Duration(rand.Intn(120)+60) * time.Second)
	}

	// if status.Status != "ok" {
	// 	panic(status)
	// }

	return nil
}

func (c *WebCLient) Like(media *instax.Media) error {
	path := fmt.Sprintf("likes/%s/like/", media.ID)

	return c.do("POST", path, media.Link)
}

func (c *WebCLient) Follow(user *instax.User) error {
	path := fmt.Sprintf("friendships/%s/follow/", user.ID)
	referer := fmt.Sprintf("https://instagram.com/%s/", user.Username)

	return c.do("POST", path, referer)
}

func main() {

	client := instax.NewClient("2079178474.1fb234f.682a311e35334df3842ccb654516baf5 ", "5ac3e50811cc47c2a4cd1adda782eb4b")
	client.Delayer = func() time.Duration {
		return time.Duration(rand.Intn(6000)+3000) * time.Millisecond
	}

	wc := NewWebClient()
	wc.CSRFToken = "266af59ac6d8c0264be518bdc4698c27"
	wc.Cookie = "mid=VIt-JAAEAAHTEAi2AXlL5hZkgvsG; ccode=RU; __utma=1.1078183635.1418427941.1432656376.1432978292.12; __utmc=1; __utmz=1.1432978292.12.2.utmcsr=t.co|utmccn=(referral)|utmcmd=referral|utmcct=/5dO7uc7mS5; sessionid=IGSCe84a66309f4b2a287b345751eff47bd0aa53f2be00fcd01e5255a64638bc4700%3Au7ZylXudt19daJSGmGJdHYiZ71nLy6s3%3A%7B%22_token_ver%22%3A1%2C%22_auth_user_id%22%3A2079178474%2C%22_token%22%3A%222079178474%3AwPP6KLe7p1XSclhb7wo7XcnFvbGhI8kI%3A5d9aa876c0335e031a4dcff2e9b054947fe23a6859d4c32a7b90a54d4acc7eb0%22%2C%22_auth_user_backend%22%3A%22accounts.backends.CaseInsensitiveModelBackend%22%2C%22last_refreshed%22%3A1441426810.129611%2C%22_platform%22%3A4%7D; ig_pr=1; ig_vw=1440; csrftoken=266af59ac6d8c0264be518bdc4698c27; ds_user_id=2079178474"
	wc.Delayer = func() time.Duration {
		return time.Duration(rand.Intn(10000)+5000) * time.Millisecond
	}

	liker := NewLiker([]string{"деньгорода"}, client, wc)
	liker.Start()

	// feed := client.MediaByTag("animallovers")
	// m, _ := feed.Get()
	// CheckAndLike(client, m)
	// for {
	// 	if count > 1 {
	// 		return
	// 	}
	// 	m, _ = feed.Next()
	// 	CheckAndLike(client, m)
	// }
}
