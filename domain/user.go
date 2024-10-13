package domain

import "time"

type User struct {
	Id       uint64    `json:"id"`
	Name     string    `json:"name"`
	Password string    `json:"password"`
	Phone    string    `json:"phone"`
	Birthday time.Time `json:"birthday"`
}
