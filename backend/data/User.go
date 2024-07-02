package data

type User struct {
	ID       int32  `json:"id"`
	Server   string `json:"server"`
	Username string `json:"username"`
	Password string `json:"-"`
}
