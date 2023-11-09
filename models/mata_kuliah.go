package models

import (
	"gorm.io/gorm"
)

type MataKuliah struct {
	gorm.Model
	Code     string `gorm:"not null" json:"code"`
	Sks      int    `gorm:"not null" json:"jumlah_sks"`
	Mk       string `gorm:"not null" json:"mata_kuliah"`
	Semester int16  `gorm:"not null" json:"semester"`
	Dosen    string `gorm:"not null" json:"dosen"`
}
