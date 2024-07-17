package models

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	ID       int
	SBD      int `gorm:"primaryKey"`
	Name     string
	Gender   string
	CMND     string
	Birthday string
	Toan     *float32
	Ly       *float32
	Hoa      *float32
	Van      *float32
	Anh      *float32
	Sinh     *float32
	Su       *float32
	Dia      *float32
	GDCD     *float32
	Trung    *float32
	Nhat     *float32
	Phap     *float32
	KHTN     *float32
	KHXH     *float32
}
