package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/012e/scrape-thptscore/stuparser"
	"github.com/glebarez/sqlite"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	ConcurrentConnection int
	StartIndex           int
	EndIndex             int
	TryNum               int
)

var f, err = os.OpenFile("logfile", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

func init() {
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(f)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	flag.IntVar(&ConcurrentConnection, "con", 1, "Number of concurrent connections")
	flag.IntVar(&StartIndex, "start", 0, "Start index")
	flag.IntVar(&EndIndex, "end", 0, "End index")
	flag.IntVar(&TryNum, "try", 3, "Total tries until give up scraping an id")
	flag.Parse()
}

func getStudent(id int) (stuparser.Student, error) {
	client := &http.Client{}
	var data = strings.NewReader(fmt.Sprintf(`layout=Decl.DataSet.Detail.default&itemsPerPage=10&pageNo=1&service=Content.Decl.DataSet.Grouping.select&itemId=64b60ae35b3cb1a42d088074&gridModuleParentId=10&type=Decl.DataSet&page=&modulePosition=0&moduleParentId=-1&orderBy=&unRegex=&keyword=%d&_t=1690715800401`, id))
	req, err := http.NewRequest("POST", "https://angiang.edu.vn/?module=Content.Listing&moduleId=1010&cmd=redraw&site=19228&url_mode=rewrite&submitFormId=1010&moduleId=1010&page=&site=19228", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/115.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Referer", "https://angiang.edu.vn/tra-cuu/diem-tot-nghiep-thpt")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Origin", "https://angiang.edu.vn")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", "be=76; AUTH_BEARER_default=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiJ9.eyJpYXQiOjE2OTA3MTI3ODksImp0aSI6IlNPbVJcL3pKZmNwMU9CM2tlMGFkcXZWRkpvXC9SR3puQ3NGNmNBTUVVdFFpaz0iLCJpc3MiOiJhbmdpYW5nLmVkdS52biIsIm5iZiI6MTY5MDcxMjc4OSwiZXhwIjoxNjkwNzE2Mzg5LCJkYXRhIjoiY3NyZlRva2VufHM6NjQ6XCJkY2Y3Y2E4YzFlNWRjYWNlNTgxYzEyNGNiNTQ4NTM2Y2IwNjZkOGVjOTkyMTUzY2QxMjg4YmE1OWU0NjgxZWZhXCI7Z3Vlc3RJZHxzOjMyOlwiMjdkNGM0ZDU3YjgyNWE5MmI3Y2RkNDYwZmMyOTA0OTZcIjt2aXNpdGVkMTkyMjh8aToxNjkwNzExOTA2OyJ9.L4xArMXuqsLEtaS3ewBqwkfvVdVNM7PQPYtf3MNkUD4tmBFkHwaXmebTjzly_g9ZbCCVn6nFKl4dryaoqnLvWQ")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return stuparser.Student{}, fmt.Errorf("failed to get response from id %d", id)
	}
	student, err := stuparser.TableFromResponse(resp)
	if err != nil {
		return stuparser.Student{}, err
	}
	return student, nil
}

func getStudentBetween(start int, end int, db *gorm.DB, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := start; i < end+1; i++ {
		log.Infof("Currently scraping student %d", i)
		var student stuparser.Student
		var err error
		for tries := int(0); tries < TryNum+1; tries++ {
			student, err = getStudent(i)
			if err != nil {
				log.Warnf("Student with id %d failed %d times", i, tries)
				time.Sleep(1 * time.Second)
			} else {
				db.Create(&student)
				break
			}
		}
		log.Infof("Finished student id %d", i)
	}
}

func efficientTaskSplit(start int, end int, con int) [][]int {
	packCount := (EndIndex - StartIndex + 1) / ConcurrentConnection
	remaining := (EndIndex - StartIndex + 1) % ConcurrentConnection
	share := make([]int, con)
	for i := 0; i < int(con); i++ {
		if remaining > 0 {
			share[i] = packCount + 1
			remaining--
		} else {
			share[i] = packCount
		}
	}
	split := make([][]int, con)
	for i := range split {
		split[i] = make([]int, 2)
	}
	for i := int(0); i < con; i++ {
		split[i][0] = start + packCount*i
		split[i][1] = split[i][0] + share[i] - 1
	}
	return split
}

func main() {
	switch {
	case ConcurrentConnection == 0:
		log.Fatal("number of concurrent connections can't be zero")
	case StartIndex == 0:
		log.Fatal("start index can't be zero")
	case EndIndex < StartIndex:
		log.Fatal("end index can't be smaller than start index")
	}
	db, err := gorm.Open(sqlite.Open("student.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// var wg sync.WaitGroup
	// wg.Add(int(ConcurrentConnection))
	// for _, split := range efficientTaskSplit(StartIndex, EndIndex, ConcurrentConnection) {
	// 	log.Infof("Starting goroutine for task between %d and %d", split[0], split[1])
	// 	go getStudentBetween(split[0], split[1], db, &wg)
	// }
	// wg.Wait()
	// defer f.Close()

	calScore(db)
}
