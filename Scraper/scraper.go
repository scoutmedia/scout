package scraper

import (
	"fmt"
	"math"
	"regexp"
	config "scout/Config"
	downloader "scout/Downloader"
	model "scout/Models"
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
	switch scraper.task.Media.Type {
	case "movie":
		if FindMovie(scraper) {
			downloader.Start(scraper.task)
		}
	}
}

func (scraper *Scraper) CreateUrl(media model.Media) string {
	replacedSpaces := strings.ReplaceAll(media.Name, " ", "%20")
	replaceCommas := strings.ReplaceAll(replacedSpaces, "'", "")
	replaceCommas += "%20" + quality

	switch media.Type {
	case "tv":
		return fmt.Sprintf("https://1337x.to/sort-category-search/%v/TV/seeders/desc/1/", replaceCommas)
	case "movie":

		return fmt.Sprintf("https://1337x.to/sort-category-search/%v/Movies/seeders/desc/1/", replaceCommas)
	}
	return ""
}

func (scraper *Scraper) GetElementData(r *colly.HTMLElement) *model.TorrentHTMLElement {
	element := &model.TorrentHTMLElement{
		Date:     r.ChildText(".coll-date"),
		Size:     r.ChildText(".size"),
		Name:     r.ChildText(".name"),
		Uploader: r.ChildText(".coll-5"),
		Hrefs:    r.ChildAttrs("a", "href"),
	}
	return element
}

func MatchName(title string, element *model.TorrentHTMLElement) bool {
	loweredTorrentName := strings.ToLower(element.Name)
	return strings.HasPrefix(loweredTorrentName, formatTitle(title, element.Uploader))
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

func ContainsNegativeWord(torrentName string, negativeWords []string) bool {
	loweredTorrentName := strings.ToLower(torrentName)
	for _, word := range negativeWords {
		if strings.Contains(loweredTorrentName, word) {
			return false
		}
		continue
	}

	return true
}

func MatchQuality(quality string, torrentName string) bool {
	loweredTorrentName := strings.ToLower(torrentName)
	loweredQuality := strings.ToLower(quality)
	return strings.Contains(loweredTorrentName, loweredQuality)
}

func MatchUploader(uploader string, sources []string) bool {
	for _, source := range sources {
		if strings.Contains(uploader, source) {
			return true
		}
		continue
	}
	return false
}

func GetFileSize(size string) (val float64, err error) {
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
