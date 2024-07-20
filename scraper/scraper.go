package scraper

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/012e/scrape-thptscore/models"
	"github.com/012e/scrape-thptscore/scrapesources/ag"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Config struct {
	DB                   *gorm.DB
	ConcurrentConnection int
	StartIndex           int
	EndIndex             int
	Retries              int
	Source               ScrapSource
}

type ScrapSource interface {
	// GetRequest returns a request to the intended scrape source
	GetRequest(sbd int) (*http.Request, error)

	// ParseResponse parses the response after getting the response from the requested source
	ParseResponse(resp *http.Response) (*models.Student, error)
}

type ErrWithID struct {
	ID  int
	Err error
}

type Scraper struct {
	inputChan  chan int
	config     Config
	errs       []ErrWithID
	errMutex   sync.Mutex
	wg         sync.WaitGroup
	httpClient *http.Client
}

func NewScraper(config Config) *Scraper {
	if config.Source == nil {
		config.Source = ag.Scraper{}
	}
	retryhttp := retryablehttp.NewClient()
	retryhttp.RetryMax = config.Retries
	client := retryhttp.StandardClient()
	return &Scraper{config: config, httpClient: client}
}

func (s *Scraper) reportError(id int, err error) {
	s.errMutex.Lock()
	s.errs = append(s.errs, ErrWithID{id, err})
	s.errMutex.Unlock()
}

func (s *Scraper) Run() {
	logrus.Debug("started scraper")
	s.inputChan = make(chan int)
	for i := 0; i < s.config.ConcurrentConnection; i += 1 {
		s.wg.Add(1)
		go s.scrape()
	}
	logrus.Debug("finished spawning goroutines")
	for i := s.config.StartIndex; i <= s.config.EndIndex; i += 1 {
		logrus.Debugf("adding %d to channel", i)
		s.inputChan <- i
	}
	logrus.Debug("closing input channel")
	close(s.inputChan)
	s.wg.Wait()
	logrus.Debug("finished scraping")
}

func (s *Scraper) getBySbd(sbd int) (*models.Student, error) {
	logrus.Debugf("getting student request for %d", sbd)
	req, err := s.config.Source.GetRequest(sbd)
	if err != nil {
		return nil, fmt.Errorf("error getting student from scraper: %v", err)
	}

	logrus.Debugf("executing request for %d", sbd)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("couldn't request student: %v", err)
	}

	logrus.Debugf("parsing request for %d", sbd)
	student, err := s.config.Source.ParseResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse student: %v", err)
	}
	return student, nil
}

func (s *Scraper) scrape() {
	for sbd := range s.inputChan {
		logrus.Debugf("Scraping student id %d", sbd)
		student, err := s.getBySbd(sbd)
		if err != nil {
			logrus.Debugf("failed to get student %d: %v", sbd, err)
			s.reportError(sbd, err)
		} else {
			logrus.Debugf("created student %v", *student)
			s.config.DB.Create(&student)
		}
	}
	s.wg.Done()
}

func (s Scraper) GetErrors() []ErrWithID {
	return s.errs
}
