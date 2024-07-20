package baoag

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/012e/scrape-thptscore/models"
	"github.com/nfx/go-htmltable"
)

type htmlTable struct {
	SBD      int    `header:"Số báo danh"`
	Name     string `header:"Tên thí sinh"`
	Birthday string `header:"Ngày sinh"`
	Gender   string `header:"Giới tính"`
	Toan     string `header:"Toán"`
	Van      string `header:"Ngữ Văn"`
	Anh      string `header:"Ngoại ngữ"`
	Ly       string `header:"Bài khoa học tự nhiên Vật lý"`
	Hoa      string `header:"Bài khoa học tự nhiên Hoá học"`
	Sinh     string `header:"Bài khoa học tự nhiên Sinh học"`
	KHTN     string `header:"Bài khoa học tự nhiên KHTN"`
	Su       string `header:"Bài khoa học xã hội Lịch sử"`
	Dia      string `header:"Bài khoa học xã hội Địa lý"`
	GDCD     string `header:"Bài khoa học xã hội GDCD"`
	KHXH     string `header:"Bài khoa học xã hội KHXH"`
}

func floatStringOrNil(s string) *float32 {
	if len(s) == 0 {
		return nil
	}
	res, err := strconv.ParseFloat(s, 32)
	if err != nil {
		panic(fmt.Sprintf("failed to parse float score: %v", err))
	}
	res32 := float32(res)
	return &res32
}

// func tableToMap(table []htmlTable) map[string]float32 {
// 	var m = make(map[string]float32)
// 	for _, v := range table {
// 		parsed, err := strconv.ParseFloat(v.Score, 32)
// 		if err != nil {
// 			panic(fmt.Sprintf("failed to parse float score: %v", err))
// 		}
// 		m[v.Subject] = float32(parsed)
// 	}
// 	return m
// }

func convSrcToDst(src htmlTable) *models.Student {
	return &models.Student{
		SBD:    src.SBD,
		Name:   src.Name,
		Gender: src.Gender,
		Toan:   floatStringOrNil(src.Toan),
		Van:    floatStringOrNil(src.Van),
		Ly:     floatStringOrNil(src.Ly),
		Hoa:    floatStringOrNil(src.Hoa),
		Sinh:   floatStringOrNil(src.Sinh),
		Su:     floatStringOrNil(src.Su),
		Dia:    floatStringOrNil(src.Dia),
		Anh:    floatStringOrNil(src.Anh),
		GDCD:   floatStringOrNil(src.GDCD),
		KHTN:   floatStringOrNil(src.KHTN),
		KHXH:   floatStringOrNil(src.KHXH),
	}
}

func parseTableFromResponse(res *http.Response) (*models.Student, error) {
	table, err := htmltable.NewSliceFromResponse[htmlTable](res)
	if err != nil {
		return nil, err
	}
	if len(table) == 0 {
		return nil, fmt.Errorf("no table found")
	}
	return convSrcToDst(table[0]), nil
}

func parseTableFromHTML(html string) (*models.Student, error) {
	table, err := htmltable.NewSliceFromString[htmlTable](html)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %v", err)
	}
	return convSrcToDst(table[0]), nil
}
