package models

type Account struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	AccessToken string `json:"-"`
}
