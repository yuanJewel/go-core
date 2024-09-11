package api

type AuthenticateConfig struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	GoogleCode string `json:"google_code"`
}

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
