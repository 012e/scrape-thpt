package ag

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/012e/scrape-thptscore/models"
	"github.com/hashicorp/go-retryablehttp"
)

type Scraper struct{}

func (s Scraper) Scrape(sbd int, retries int) (*models.Student, error) {
	response, err := s.getStudentResp(sbd, retries)
	if err != nil {
		return nil, fmt.Errorf("failed to get response from id %d: %e", sbd, err)
	}
	student, err := TableFromResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response from id %d: %e", sbd, err)
	}
	student.SBD = sbd
	return student, nil
}

func (s Scraper) getStudentResp(sbd int, retries int) (*http.Response, error) {
	retrier := retryablehttp.NewClient()
	retrier.RetryMax = retries
	client := retrier.StandardClient()
	var data = strings.NewReader(fmt.Sprintf(`layout=Decl.DataSet.Detail.default&itemsPerPage=10&pageNo=1&service=Content.Decl.DataSet.Grouping.select&itemId=64b60ae35b3cb1a42d088074&gridModuleParentId=10&type=Decl.DataSet&page=&modulePosition=0&moduleParentId=-1&orderBy=&unRegex=&keyword=%d&_t=1721134551628`, sbd))
	req, err := http.NewRequest("POST", "https://angiang.edu.vn/?module=Content.Listing&moduleId=1010&cmd=redraw&site=19228&url_mode=rewrite&submitFormId=1010&moduleId=1010&page=&site=19228", data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Referer", "https://angiang.edu.vn/tra-cuu/diem-tot-nghiep-thpt")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Origin", "https://angiang.edu.vn")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", "AUTH_BEARER_default=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiJ9.eyJpYXQiOjE3MjExMzM1NzAsImp0aSI6ImNoVjFsT3RcL3lJNUh0U3NXUG1BK2tlMXl0TGRcLytvSHhMaWdaSDFmZm50cz0iLCJpc3MiOiJhbmdpYW5nLmVkdS52biIsIm5iZiI6MTcyMTEzMzU3MCwiZXhwIjoxNzIxMTM3MTcwLCJkYXRhIjoiY3NyZlRva2VufHM6NjQ6XCIwZjdkMDMxMTY5OWFiOGRlNmU1Y2U2NjlkZWYxNmVmOGQ5YWQxODNkNDY5ZTBiMzBlNzljZTVmZWRmZDc2NTZlXCI7Z3Vlc3RJZHxzOjMyOlwiMTEwNjAyYjE2Njg5YTEwOWE4MTJmNWY1MGU1ODlkZDZcIjt2aXNpdGVkMTkyMjh8aToxNzIxMTMzNTQ3OyJ9.ZQNKpODQQs4dJ3NiTexLNZKZk0MTu6otnBUakbWQoWZ5mEkA_C2ili9x8imbw5FFq5RvV6UVrgL3TdGvw3Vx8g")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Priority", "u=1")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
