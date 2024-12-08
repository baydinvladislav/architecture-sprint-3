package web

type NewUserIn struct {
	Username string
	Password string
}

type NewUserOut struct {
	//ID           uint `json:"id"`
	Username     string `json:"username"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
