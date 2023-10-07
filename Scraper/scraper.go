package scraper

import (
	"errors"
	"fmt"
	"log"
	"math"
	"regexp"
	config "scout/Config"
	downloader "scout/Downloader"
	"scout/Models"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

var (
	quality  string = "1080p"
	username string = "TGxGoodies"
)

type Scraper struct {
	c       *colly.Collector
	sources []string
	File    Models.TorrentFile
}

func NewScraper(collector *colly.Collector, sources []string) *Scraper {
	return &Scraper{
		c:       collector,
		sources: sources,
		File:    Models.TorrentFile{},
	}
}

func (scraper *Scraper) Init(requestedTitle string) {
	scrapedFile, err := scraper.FindRequestedFile(requestedTitle)
	if err != nil {
		log.Println(err)
		return
	}
	config := config.NewConfig()
	config.Load()
	d := downloader.NewDownloader(config.DataDir)
	d.Init(requestedTitle, scrapedFile)
}

func (scraper *Scraper) FindRequestedFile(title string) (file Models.TorrentFile, nil error) {
	scraper.c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	})

	scraper.c.OnHTML("table", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, r *colly.HTMLElement) {
			scraper.verifyTorrentOption(title, r)
		})
	})
	scraper.c.Visit(formatUrl(title))
	if scraper.File.Name == "" {
		return scraper.File, errors.New("Requested title is unavailable")
	}
	return scraper.File, nil
}

func formatUrl(title string) string {
	replacedSpaces := strings.ReplaceAll(title, " ", "%20")
	replaceCommas := strings.ReplaceAll(replacedSpaces, "'", "")
	return fmt.Sprintf("https://1337x.to/sort-category-search/%v/Movies/seeders/desc/1/", replaceCommas)
}

func (scraper *Scraper) verifyTorrentOption(title string, r *colly.HTMLElement) {
	uploader := r.ChildText(".coll-5")
	fileName := r.ChildText(".name")
	date := r.ChildText(".coll-date")
	size := r.ChildText(".size")
	hrefs := r.ChildAttrs("a", "href")

	if hasUsername(uploader, scraper.sources) {
		if hasQuality(quality, fileName) {
			if hasMatchingTitle(title, fileName, scraper.sources) {
				log.Println(fileName)
				if scraper.File.Name == "" {
					scraper.File.Name = fileName
					scraper.File.Date = date
					scraper.File.Uploader = username
					size, err := getFileSize(size)
					if err != nil {
						log.Println(err)
					}
					scraper.File.Size = size
					scraper.c.OnHTML("a", func(link *colly.HTMLElement) {
						if link.Text == "Magnet Download" {
							scraper.File.Magnet = link.Attr("href")
						}
					})
					scraper.c.Visit("https://1337x.to" + hrefs[1])
				} else {
					size, err := getFileSize(size)
					if err != nil {
						log.Println(err)
					}
					if size < scraper.File.Size {
						scraper.File.Name = fileName
						scraper.File.Date = date
						scraper.File.Uploader = username
						scraper.File.Size = size
						scraper.c.OnHTML("a", func(link *colly.HTMLElement) {
							if link.Text == "Magnet Download" {
								scraper.File.Magnet = link.Attr("href")
							}
						})
						scraper.c.Visit("https://1337x.to" + hrefs[1])
					}
				}
			}
			scraper.c.Visit("https://1337x.to" + hrefs[1])
		}
	}
}

func hasUsername(uploader string, sources []string) bool {
	for _, source := range sources {
		if strings.Contains(uploader, source) {
			return true
		}
		continue
	}
	return false
}

func hasMatchingTitle(title string, fileName string, sources []string) bool {
	for _, source := range sources {
		configuredTitle := configureTitle(title, source)
		if strings.Contains(fileName, configuredTitle) {
			return true
		}
	}
	return false
}

func configureTitle(title string, source string) string {
	switch source {
	case "TGxGoodies":
		r := regexp.MustCompile("\\(|\\)")
		r2 := regexp.MustCompile(" ")
		r3 := regexp.MustCompile("\\'")
		replacedStr := r3.ReplaceAllString(r2.ReplaceAllString(r.ReplaceAllString(title, ""), "."), "")
		return replacedStr
	}

	return ""
}

func hasQuality(quality string, fileName string) bool {
	if strings.Contains(fileName, quality) {
		return true
	}
	return false
}

func getFileSize(size string) (val float64, err error) {
	if strings.Contains(size, "GB") {
		trimmedStr := strings.Split(size, " ")
		if s, err := strconv.ParseFloat(trimmedStr[0], 32); err == nil {
			fixedVal := toFixed(s, 2)
			return fixedVal, err
		}
	}
	return val, err
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return math.Round(num*output) / output
}
