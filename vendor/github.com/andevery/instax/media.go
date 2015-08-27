package instax

type Media struct {
	Comments struct {
		Count int `json:"count"`
		Data  []struct {
			CreatedTime string `json:"created_time"`
			From        struct {
				FullName       string `json:"full_name"`
				ID             string `json:"id"`
				ProfilePicture string `json:"profile_picture"`
				Username       string `json:"username"`
			} `json:"from"`
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"data"`
	} `json:"comments"`
	ID    string `json:"id"`
	Likes struct {
		Count int `json:"count"`
		Data  []struct {
			FullName       string `json:"full_name"`
			ID             string `json:"id"`
			ProfilePicture string `json:"profile_picture"`
			Username       string `json:"username"`
		} `json:"data"`
	} `json:"likes"`
	Tags []string `json:"tags"`
	Type string   `json:"type"`
	User struct {
		FullName       string `json:"full_name"`
		ID             string `json:"id"`
		ProfilePicture string `json:"profile_picture"`
		Username       string `json:"username"`
	} `json:"user"`
	UserHasLiked bool `json:"user_has_liked"`
}
