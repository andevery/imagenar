package main

import (
	"github.com/andevery/instax"
	"log"
)

func main() {
	// var m instax.Response

	// client := new(http.Client)

	// requestUrl := new(url.URL)
	// requestUrl.Scheme = "https"
	// requestUrl.Host = "api.instagram.com/v1"
	// requestUrl.Path = "media/1059995724270397093_2646144"
	// q := requestUrl.Query()
	// q.Add("access_token", "2079178474.1fb234f.682a311e35334df3842ccb654516baf5")
	// requestUrl.RawQuery = q.Encode()
	// log.Println(requestUrl.String())
	// req, err := http.NewRequest("GET", requestUrl.String(), nil)
	// if err != nil {
	// 	log.Println(err)
	// }

	// resp, err := client.Do(req)
	// // resp, err := http.Get("https://api.instagram.com/v1/media/1059995724270397093_2646144?access_token=2079178474.1fb234f.682a311e35334df3842ccb654516baf5")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// body, err := ioutil.ReadAll(resp.Body)
	// err = json.Unmarshal(body, &m)
	// if err != nil {
	// 	log.Println(err)
	// }

	// log.Println(m)

	// // var media instax.Media
	// // if m.Meta.Code == 200 {
	// // 	err = json.Unmarshal(m.Data, &media)
	// // 	if err != nil {
	// // 		log.Println(err)
	// // 	}
	// // 	fmt.Println(media.Comments.Data[0].Text)
	// // }
	// log.Println(m)
	// insta := instax.NewClient("2079178474.1fb234f.682a311e35334df3842ccb654516baf5")
	// log.Println(insta)

	client := instax.NewClient("2079178474.1fb234f.682a311e35334df3842ccb654516baf5")
	feed := client.MediaByTag("moscow")

	_, err := feed.Get()
	if err != nil {
		log.Println(err)
	}
	log.Println(feed)
	log.Println(client.Limit())

	_, err = feed.Next()
	if err != nil {
		log.Println(err)
	}
	log.Println(feed)
	log.Println(client.Limit())

	feed = client.MediaForUser("self")
	m, err := feed.Get()
	if err != nil {
		log.Println(err)
	}

	err = client.Like(m[0].ID)
	if err != nil {
		log.Println(err)
	}

	log.Println(client.Limit())
}
