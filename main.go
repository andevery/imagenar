package main

import (
	"github.com/andevery/autogram"
	"github.com/andevery/instaw"
	"github.com/andevery/instax"
	"log"
	"time"
)

func main() {

	// client := instax.NewClient("2079178474.1fb234f.682a311e35334df3842ccb654516baf5", "5ac3e50811cc47c2a4cd1adda782eb4b")
	// client.Delayer = func() time.Duration {
	// 	return time.Duration(rand.Intn(6000)+3000) * time.Millisecond
	// }

	// wc := instaw.NewClient()
	// wc.CSRFToken = "266af59ac6d8c0264be518bdc4698c27"
	// wc.Cookie = "mid=VIt-JAAEAAHTEAi2AXlL5hZkgvsG; ccode=RU; __utma=1.1078183635.1418427941.1432656376.1432978292.12; __utmc=1; __utmz=1.1432978292.12.2.utmcsr=t.co|utmccn=(referral)|utmcmd=referral|utmcct=/5dO7uc7mS5; sessionid=IGSCe84a66309f4b2a287b345751eff47bd0aa53f2be00fcd01e5255a64638bc4700%3Au7ZylXudt19daJSGmGJdHYiZ71nLy6s3%3A%7B%22_token_ver%22%3A1%2C%22_auth_user_id%22%3A2079178474%2C%22_token%22%3A%222079178474%3AwPP6KLe7p1XSclhb7wo7XcnFvbGhI8kI%3A5d9aa876c0335e031a4dcff2e9b054947fe23a6859d4c32a7b90a54d4acc7eb0%22%2C%22_auth_user_backend%22%3A%22accounts.backends.CaseInsensitiveModelBackend%22%2C%22last_refreshed%22%3A1441426810.129611%2C%22_platform%22%3A4%7D; ig_pr=1; ig_vw=1440; csrftoken=266af59ac6d8c0264be518bdc4698c27; ds_user_id=2079178474"
	// wc.Delayer = func() time.Duration {
	// 	return time.Duration(rand.Intn(10000)+5000) * time.Millisecond
	// }

	// liker := NewLiker([]string{"деньгорода"}, client, wc)
	// liker.Start()

	api := instax.NewClient("2079178474.1fb234f.682a311e35334df3842ccb654516baf5")
	web, err := instaw.NewClient("andy_odds", "Dont_Panic1")
	if err != nil {
		log.Fatal(err)
	}

	l := autogram.DefaultLimiter(api, web)

	fp, err := autogram.NewProvider(autogram.MEDIA, []string{"283175195"}, l)
	if err != nil {
		log.Fatal(err)
	}

	lp, err := autogram.EmptyProvider(l)
	if err != nil {
		log.Fatal(err)
	}

	liker := autogram.DefaultLiker(lp)
	follower := autogram.DefaultFollower(fp, liker)

	// us, err := apiClient.Likes("1063050685134653313")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	go follower.Start()
	for {
		log.Printf("Likes: %v\tFollows: %v\n", liker.Provider.TotalAmount(), follower.Provider.TotalAmount())
		time.Sleep(30 * time.Second)
	}
}
