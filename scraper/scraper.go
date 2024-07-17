package scraper

import (
	"sync"

	"github.com/012e/scrape-thptscore/models"
	"github.com/012e/scrape-thptscore/scrapesources/ag"
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
	Scrape(sbd int, tries int) (*models.Student, error)
}

type ErrWithID struct {
	ID  int
	Err error
}

type Scraper struct {
	inputChan chan int
	config    Config
	errs      []ErrWithID
	errMutex  sync.Mutex
	wg        sync.WaitGroup
}

func NewScraper(config Config) *Scraper {
	if config.Source == nil {
		config.Source = ag.Scraper{}
	}
	return &Scraper{config: config}
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
		go s.Scrape()
	}
	logrus.Debug("finished spawning goroutines")
	for i := s.config.StartIndex; i <= s.config.EndIndex; i += 1 {
		logrus.Debugf("adding %d to channel", i)
		s.inputChan <- i
	}
	logrus.Debug("closing input channel")
	close(s.inputChan)
	s.wg.Wait()
}

func (s *Scraper) Scrape() {
	for sbd := range s.inputChan {
		logrus.Debugf("Scraping student id %d", sbd)
		student, err := s.config.Source.Scrape(sbd, s.config.Retries)
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

func (s *Scraper) GetErrors() []ErrWithID {
	return s.errs
}
