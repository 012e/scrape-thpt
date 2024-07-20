package vietnamnet

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/012e/scrape-thptscore/models"
	"github.com/nfx/go-htmltable"
)

type htmlTable struct {
	Subject string `header:"Môn"`
	Score   string `header:"Điểm"`
}

func okOrNil(m map[string]float32, key string) *float32 {
	if v, ok := m[key]; ok {
		return &v
	}
	return nil
}

func tableToMap(table []htmlTable) map[string]float32 {
	var m = make(map[string]float32)
	for _, v := range table {
		parsed, err := strconv.ParseFloat(v.Score, 32)
		if err != nil {
			panic(fmt.Sprintf("failed to parse float score: %v", err))
		}
		m[v.Subject] = float32(parsed)
	}
	return m
}

func (s Scraper) ParseResponse(res *http.Response) (*models.Student, error) {
	table, err := htmltable.NewSliceFromResponse[htmlTable](res)
	if err != nil {
		return nil, err
	}
	return mapToStudent(tableToMap(table))
}

func parseTableFromHTML(html string) (*models.Student, error) {
	table, err := htmltable.NewSliceFromString[htmlTable](html)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %v", err)
	}
	return mapToStudent(tableToMap(table))
}

func mapToStudent(scores map[string]float32) (*models.Student, error) {
	return &models.Student{
		Toan:  okOrNil(scores, "Toán"),
		Van:   okOrNil(scores, "Văn"),
		Ly:    okOrNil(scores, "Lí"),
		Hoa:   okOrNil(scores, "Hóa"),
		Sinh:  okOrNil(scores, "Sinh"),
		Su:    okOrNil(scores, "Sử"),
		Dia:   okOrNil(scores, "Địa"),
		Anh:   okOrNil(scores, "Ngoại ngữ"),
		GDCD:  okOrNil(scores, "GDCD"),
		Nhat:  okOrNil(scores, "Tiếng Nhật"),
		Phap:  okOrNil(scores, "Tiếng Pháp"),
		Trung: okOrNil(scores, "Tiếng Trung"),
		KHTN:  okOrNil(scores, "KHTN"),
		KHXH:  okOrNil(scores, "KHXH"),
	}, nil
}
