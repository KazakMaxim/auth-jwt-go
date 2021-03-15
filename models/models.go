package models

type User struct {
	User_guid string
	Username  string `json:"username" binding "required"`
	Password  string `json:"password" binding "required"`
}

type RefreshTable struct {
	User_guid     string `json:"user_guid" binding "required"`
	Refresh_token string `json:"refresh_token" binding "required"`
}

type Tokens struct {
	Access  string `json:"access_token" binding "required"`
	Refresh string `json:"refresh_token" binding "required"`
}
