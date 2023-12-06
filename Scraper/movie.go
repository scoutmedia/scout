package scraper

import (
	"log"

	"github.com/gocolly/colly/v2"
)

func FindMovie(scraper *Scraper) bool {
	scraper.c.OnHTML("table", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, r *colly.HTMLElement) {
			checkTorrent(scraper, r)
		})
	})

	err := scraper.c.Visit(scraper.CreateUrl(scraper.task.Media))
	if err != nil {
		log.Println(err)
	}
	return scraper.task.TorrentFile.Magnet != ""
}

func checkTorrent(scraper *Scraper, r *colly.HTMLElement) {
	element := scraper.GetElementData(r)

	if MatchUploader(element.Uploader, scraper.task.Sources) {
		if ContainsNegativeWord(element.Name, scraper.config.NegativeWords) {
			if MatchQuality(scraper.config.Quality, element.Name) {
				if MatchName(scraper.task.Media.Name, element) {
					if scraper.task.TorrentFile.Name == "" {
						scraper.task.TorrentFile.Name = element.Name
						scraper.task.TorrentFile.Date = element.Date
						scraper.task.TorrentFile.Uploader = element.Uploader
						size, err := GetFileSize(element.Size)
						if err != nil {
							log.Println("Error getting file size", err, element.Name)
						}
						scraper.task.TorrentFile.Size = size
						scraper.c.OnHTML("a", func(link *colly.HTMLElement) {
							if link.Text == "Magnet Download" {
								scraper.task.TorrentFile.Magnet = link.Attr("href")
							}
						})
						scraper.c.Visit("https://1337x.to" + element.Hrefs[1])
					}
					size, err := GetFileSize(element.Size)
					if err != nil {
						log.Println(err)
					}
					if size < scraper.task.TorrentFile.Size {
						scraper.task.TorrentFile.Name = element.Name
						scraper.task.TorrentFile.Date = element.Date
						scraper.task.TorrentFile.Uploader = element.Uploader
						scraper.task.TorrentFile.Size = size
						scraper.c.OnHTML("a", func(link *colly.HTMLElement) {
							if link.Text == "Magnet Download" {
								scraper.task.TorrentFile.Magnet = link.Attr("href")
							}
						})
						scraper.c.Visit("https://1337x.to" + element.Hrefs[1])
					}
				}
			}
		}
	}
}
