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
	if scraper.task.TorrentFile.Name != "" {
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
		if matchQuality(quality, torrentName) {
			log.Println(scraper.task.Name)
			log.Println(matchTitle(scraper.task.Name, torrentName, uploader))
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
					}
					scraper.c.Visit("https://1337x.to" + hrefs[1])
				}

			}
		}
	}
}

func checkFileSize(torrentSize string, currentFileSize float64) bool {
	convertedFileSize, err := getFileSize(torrentSize)
	if err != nil {
		log.Println(err)
	}
	if convertedFileSize < currentFileSize {
		return true
	}
	return false
}

func matchTitle(title string, torrentName string, uploader string) bool {
	if strings.Contains(torrentName, formatTitle(title, uploader)) {
		return true
	}
	return false
}

func formatTitle(title string, uploader string) (replacedStr string) {
	switch uploader {
	case "TGxGoodies":
		r := regexp.MustCompile("\\(|\\)")
		r2 := regexp.MustCompile(" ")
		r3 := regexp.MustCompile("\\'")
		r4 := regexp.MustCompile("\\:")
		replacedStr = r4.ReplaceAllString(r3.ReplaceAllString(r2.ReplaceAllString(r.ReplaceAllString(title, ""), "."), ""), "")
		log.Println(replacedStr)
		return replacedStr
	}
	return replacedStr
}

func matchQuality(quality string, torrentName string) bool {
	loweredTorrentName := strings.ToLower(torrentName)
	loweredQuality := strings.ToLower(quality)

	if strings.Contains(loweredTorrentName, loweredQuality) {
		return true
	}
	return false
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

// func (scraper *Scraper) FindRequestedFile(title string) (file Models.TorrentFile, nil error) {
// 	scraper.c.OnRequest(func(r *colly.Request) {
// 		r.Headers.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
// 	})

// 	scraper.c.OnHTML("table", func(e *colly.HTMLElement) {
// 		e.ForEach("tr", func(i int, r *colly.HTMLElement) {
// 			scraper.verifyTorrentOption(title, r)
// 		})
// 	})
// 	scraper.c.Visit(formatUrl(title))
// 	if scraper.File.Name == "" {
// 		return scraper.File, errors.New("Requested title is unavailable")
// 	}
// 	return scraper.File, nil
// }

// func (scraper *Scraper) verifyTorrentOption(title string, r *colly.HTMLElement) {
// 	uploader := r.ChildText(".coll-5")
// 	fileName := r.ChildText(".name")
// 	date := r.ChildText(".coll-date")
// 	size := r.ChildText(".size")
// 	hrefs := r.ChildAttrs("a", "href")

// 	if hasUsername(uploader, scraper.sources) {
// 		if hasQuality(quality, fileName) {
// 			if hasMatchingTitle(title, fileName, scraper.sources) {
// 				if scraper.File.Name == "" {
// 					scraper.File.Name = fileName
// 					scraper.File.Date = date
// 					scraper.File.Uploader = username
// 					size, err := getFileSize(size)
// 					if err != nil {
// 						log.Println(err)
// 					}
// 					scraper.File.Size = size
// 					scraper.c.OnHTML("a", func(link *colly.HTMLElement) {
// 						if link.Text == "Magnet Download" {
// 							scraper.File.Magnet = link.Attr("href")
// 						}
// 						})
// 					scraper.c.Visit("https://1337x.to" + hrefs[1])
// 				} else {
// 					size, err := getFileSize(size)
// 					if err != nil {
// 						log.Println(err)
// 					}
// 					if size < scraper.File.Size {
// 						scraper.File.Name = fileName
// 						scraper.File.Date = date
// 						scraper.File.Uploader = username
// 						scraper.File.Size = size
// 						scraper.c.OnHTML("a", func(link *colly.HTMLElement) {
// 							if link.Text == "Magnet Download" {
// 								scraper.File.Magnet = link.Attr("href")
// 							}
// 						})
// 						scraper.c.Visit("https://1337x.to" + hrefs[1])
// 					}
// 				}
// 				// scraper.c.Visit("https://1337x.to" + hrefs[1])
// 			}
// 		}
// 	}
// }

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
