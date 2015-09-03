package instax

type User struct {
	Bio    string `json:"bio"`
	Counts struct {
		FollowedBy int `json:"followed_by"`
		Follows    int `json:"follows"`
		Media      int `json:"media"`
	} `json:"counts"`
	FullName       string `json:"full_name"`
	ID             string `json:"id"`
	ProfilePicture string `json:"profile_picture"`
	Username       string `json:"username"`
	Website        string `json:"website"`
}
