package baoag

import (
	"fmt"
	"net/http"
)

type Scraper struct{}

func (_ Scraper) GetRequest(sbd int) (*http.Request, error) {
	return http.NewRequest(
		"GET",
		fmt.Sprintf("https://baoangiang.com.vn/tra-cuu-diem-thi-thpt.html?tensbd=%d&cumthi=undefined", sbd),
		nil,
	)
}
