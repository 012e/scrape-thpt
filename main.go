package main

import (
	"flag"
	"os"

	"github.com/012e/scrape-thptscore/models"
	"github.com/012e/scrape-thptscore/scraper"
	"github.com/012e/scrape-thptscore/scrapesources/vietnamnet"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ConcurrentConnection int
	StartIndex           int
	EndIndex             int
	Retries              int
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	flag.IntVar(&ConcurrentConnection, "con", 1, "Number of concurrent connections")
	flag.IntVar(&StartIndex, "start", 0, "Start index")
	flag.IntVar(&EndIndex, "end", 0, "End index")
	flag.IntVar(&Retries, "try", 3, "Total tries until give up scraping an id")
	flag.Parse()
	if StartIndex <= 0 {
		log.Fatal("Invalid start index (not specified or less than 0)")
	}
	if EndIndex < StartIndex {
		EndIndex = StartIndex
	}
}

func runMigration(db *gorm.DB) {
	db.AutoMigrate(&models.Student{})
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

	log.Debug("connecting to db")
	db, err := gorm.Open(sqlite.Open("student.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Debug("connected to db")

	runMigration(db)
	log.Debug("finished db migrations")

	scraper := scraper.NewScraper(scraper.Config{
		ConcurrentConnection: ConcurrentConnection,
		StartIndex:           StartIndex,
		EndIndex:             EndIndex,
		Retries:              Retries,
		DB:                   db,
		Source:               vietnamnet.Scraper{},
	})
	scraper.Run()
	errs := scraper.GetErrors()
	for _, e := range errs {
		log.Errorf("failed to scrape id %d: %v", e.ID, e.Err)
	}
}
