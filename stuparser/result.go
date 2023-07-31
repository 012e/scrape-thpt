package stuparser

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	ID             int    `json:"id" header:"Số báo danh"`
	Name           string `json:"name" header:"Tên"`
	UnparsedResult string `json:"unparsedResult" header:"Điểm chi tiết"`
	BirthdayStr    string `json:"birthdayStr" header:"Ngày sinh"`
	Gender         string `json:"gender" header:"Giới tính"`
	CMND           string `json:"cmnd" header:"Số CMND"`
	Toan           float32
	Ly             float32
	Hoa            float32
	Van            float32
	Anh            float32
	Sinh           float32
	Su             float32
	Dia            float32
	GDCD           float32
	Trung          float32
	Nhat           float32
	Phap           float32
	KHTN           float32
	KHXH           float32
}
