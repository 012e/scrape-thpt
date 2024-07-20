package vietnamnet

import (
	"fmt"
	"net/http"
)

type Scraper struct{}

func (s Scraper) GetRequest(sbd int) (*http.Request, error) {
	return http.NewRequest(http.MethodGet,
		fmt.Sprintf("https://vietnamnet.vn/giao-duc/diem-thi/tra-cuu-diem-thi-tot-nghiep-thpt/2024/%d.html", sbd),
		nil)
}
