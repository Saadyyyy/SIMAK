package models

import "gorm.io/gorm"

type Person struct {
	gorm.Model
	Nama        string `json:"nama"`
	Nim         string `json:"nim"`
	Password    string `json:"password"`
	Prodi       string `json:"prodi"`
	PhoneNumber string `json:"phone_number"`
	Semester    int16  `json:"semster"`
	IsAdmin     bool   `json:"is_admin"`
}
