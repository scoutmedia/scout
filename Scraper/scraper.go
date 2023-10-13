package scraper

import (
	"fmt"
	"log"
	"math"
	"regexp"
	config "scout/Config"
	downloader "scout/Downloader"
	task "scout/Task"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

var (
	quality string = "1080p"
)

type Scraper struct {
	c      *colly.Collector
	task   *task.Task
	config *config.Config
}

func NewScraper(collector *colly.Collector, task *task.Task, config config.Config) *Scraper {
	collector.UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36"
	return &Scraper{
		c:      collector,
		task:   task,
		config: &config,
	}
}

func (scraper *Scraper) Start(downloader *downloader.Downloader) {
	if scraper.find() {
		downloader.Start(scraper.task.Name, scraper.task.TorrentFile)
	}
}

func (scraper *Scraper) find() bool {
	scraper.c.OnHTML("table", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, r *colly.HTMLElement) {
			scraper.verify(r)
		})
	})
	scraper.c.Visit(createsearchUrl(scraper.task.Name))
	if scraper.task.Name != "" {
		return true
	}
	return false
}

func (scraper *Scraper) verify(r *colly.HTMLElement) {
	uploader := r.ChildText(".coll-5")
	torrentName := r.ChildText(".name")
	uploadDate := r.ChildText(".coll-date")
	torrentSize := r.ChildText(".size")
	hrefs := r.ChildAttrs("a", "href")

	if matchUsername(uploader, scraper.task.Sources) {
		if containsNegativeWord(torrentName, scraper.config.NegativeWords) {
			if matchQuality(quality, torrentName) {
				if matchTitle(scraper.task.Name, torrentName, uploader) {
					if scraper.task.TorrentFile.Name == "" {
						scraper.task.TorrentFile.Name = torrentName
						scraper.task.TorrentFile.Date = uploadDate
						scraper.task.TorrentFile.Uploader = uploader
						size, err := getFileSize(torrentSize)
						if err != nil {
							log.Println(err)
						}
						scraper.task.TorrentFile.Size = size
						scraper.c.OnHTML("a", func(link *colly.HTMLElement) {
							if link.Text == "Magnet Download" {
								scraper.task.TorrentFile.Magnet = link.Attr("href")
							}
						})
						scraper.c.Visit("https://1337x.to" + hrefs[1])
					} else {
						size, err := getFileSize(torrentSize)
						if err != nil {
							log.Println(err)
						}
						if size < scraper.task.TorrentFile.Size {
							scraper.task.TorrentFile.Name = torrentName
							scraper.task.TorrentFile.Date = uploadDate
							scraper.task.TorrentFile.Uploader = uploader
							scraper.task.TorrentFile.Size = size
							scraper.c.OnHTML("a", func(link *colly.HTMLElement) {
								if link.Text == "Magnet Download" {
									scraper.task.TorrentFile.Magnet = link.Attr("href")
								}
							})
							scraper.c.Visit("https://1337x.to" + hrefs[1])
						}
					}

				}
			}
		}
	}
}

func matchTitle(title string, torrentName string, uploader string) bool {
	loweredTorrentName := strings.ToLower(torrentName)
	return strings.HasPrefix(loweredTorrentName, formatTitle(title, uploader))
}

func formatTitle(title string, uploader string) (replacedStr string) {
	switch uploader {
	case "TGxGoodies":
		r := regexp.MustCompile("\\(|\\)")
		r2 := regexp.MustCompile(" ")
		r3 := regexp.MustCompile("\\'")
		r4 := regexp.MustCompile("\\:")
		replacedStr = r4.ReplaceAllString(r3.ReplaceAllString(r2.ReplaceAllString(r.ReplaceAllString(title, ""), "."), ""), "")
		return strings.ToLower(replacedStr)
	}
	return replacedStr
}

func containsNegativeWord(torrentName string, negativeWords []string) bool {
	loweredTorrentName := strings.ToLower(torrentName)
	for _, word := range negativeWords {
		if strings.Contains(loweredTorrentName, word) {
			return false
		}
		continue
	}

	return true
}

func matchQuality(quality string, torrentName string) bool {
	loweredTorrentName := strings.ToLower(torrentName)
	loweredQuality := strings.ToLower(quality)
	return strings.Contains(loweredTorrentName, loweredQuality)
}

func matchUsername(uploader string, sources []string) bool {
	for _, source := range sources {
		if strings.Contains(uploader, source) {
			return true
		}
		continue
	}
	return false
}

func createsearchUrl(title string) string {
	replacedSpaces := strings.ReplaceAll(title, " ", "%20")
	replaceCommas := strings.ReplaceAll(replacedSpaces, "'", "")
	replaceCommas += "%20" + quality

	return fmt.Sprintf("https://1337x.to/sort-category-search/%v/Movies/seeders/desc/1/", replaceCommas)
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
