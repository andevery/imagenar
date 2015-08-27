package main

import (
	"encoding/json"
	"fmt"
	"github.com/andevery/instax"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func main() {
	var m instax.Response

	client := new(http.Client)

	requestUrl := new(url.URL)
	requestUrl.Scheme = "https"
	requestUrl.Host = "api.instagram.com"
	requestUrl.Path = "v1/media/1059995724270397093_2646144"
	requestUrl.Query().Add("access_token", "2079178474.1fb234f.682a311e35334df3842ccb654516baf5")
	req, err := http.NewRequest("GET", requestUrl.String(), nil)
	if err != nil {
		log.Println(err)
	}
	log.Println(req.URL)

	resp, err := client.Do(req)
	// resp, err := http.Get("https://api.instagram.com/v1/media/1059995724270397093_2646144?access_token=2079178474.1fb234f.682a311e35334df3842ccb654516baf5")
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(m)
}
