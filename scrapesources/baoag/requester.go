package baoag

import (
	"fmt"
	"net/http"

	"github.com/012e/scrape-thptscore/models"
	"github.com/hashicorp/go-retryablehttp"
)

type Scraper struct{}

func getResponse(sbd int, retries int) (*http.Response, error) {
	client := retryablehttp.NewClient()
	client.RetryMax = retries
	return client.Get(fmt.Sprintf("https://baoangiang.com.vn/tra-cuu-diem-thi-thpt.html?tensbd=%d&cumthi=undefined", sbd))
}

func (s Scraper) Scrape(sdb int, retries int) (*models.Student, error) {
	resp, err := getResponse(sdb, retries)
	if err != nil {
		return nil, fmt.Errorf("failed to get response from id %d: %e", sdb, err)
	}
	student, err := parseTableFromResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response from id %d: %e", sdb, err)
	}
	student.SBD = sdb
	return student, nil
}
