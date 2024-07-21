# Scrape THPTQG score

Scrape scores from Vietnam's national high school exam with any source of choice.
![image](https://github.com/user-attachments/assets/f0833c79-d1fa-4e06-a255-1e2da9c5b68c)

## Installation

### As a library

Simply run:

```sh
go get github.com/012e/scrape-thpt
```

### As a binary

```sh
go install github.com/012e/scrape-thpt@latest
```

## Usage

### As a binary

By default, scrapes score from [Báo An Giang](https://baoangiang.com.vn/tra-cuu-diem-thi-thpt.html)
and saves them to `students` table.
Basic usage:
```txt
$ scrape-thpt -help
  -con int
        Number of concurrent connections.
        Tweaks this number to scrape faster. (default 3)
  -end int
        End index, default value is start index
  -start int
        Start index
  -try int
        Total tries until give up scraping a candidate number (default 3)
```

### As a library

#### `scrapesource` interface

All the scraping is mainly based on `ScrapeSource` interface:
```go
type ScrapSource interface {
    // GetRequest returns a request to the intended scrape source
    GetRequest(sbd int) (*http.Request, error)

    // ParseResponse parses the response after getting the response from the requested source
    ParseResponse(resp *http.Response) (*models.Student, error)
}
```

For examples checkout [scrapesources](./scrapesources). It has already implemented
scrape sources from [baoangiang.com.vn](https://baoangiang.com.vn/tra-cuu-diem-thi-thpt.html),
[vietnamnet.vn](https://vietnamnet.vn/giao-duc/diem-thi/tra-cuu-diem-thi-tot-nghiep-thpt-2023),
[angiang.edu.vn](https://angiang.edu.vn/tra-cuu/diem-tot-nghiep-thpt).

#### Scraping

1. To begin, create a new `Scraper`:
```go
db := createGormDB()
scraper := scraper.NewScraper(scraper.Config{
    ConcurrentConnection: 3,               // three goroutines for scraping
    StartIndex:           51000001,        // scrapes between those range
    EndIndex:             51000010,
    Retries:              3,         // will retries 3 times before giving up
    DB:                   db,              // any gorm instance
    Source:               baoag.Scraper{}, // anything implements `ScrapeSource` interface
})
```
Currently, gorm is the only supported orm.

2. Start scraping and handle errors (slice of structs contain error and candidate number):
```go
scraper.Run()
errors := scraper.GetErrors()
for _, err := range errors {
	fmt.Println("failed %d: %v", err.ID, err.Err)
	// handle error
}
```


## License

MIT license.
