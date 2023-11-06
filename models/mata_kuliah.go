package models

import "gorm.io/gorm"

type MataKuliah struct {
	gorm.Model
	Mk       string `json:"mata_kuliah"`
	Semester int16  `json:"semester"`
	Dosen    string `json:"dosen"`
}
