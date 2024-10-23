package domain

import "time"

type User struct {
	Id         uint64
	Name       string
	Password   string
	Phone      string
	IsMerchant bool
	Birthday   time.Time
}

type Address struct {
	Id        uint64
	UserId    uint64
	Street    string
	City      string
	State     string
	ZipCode   string
	Country   string
	IsDefault bool
}
