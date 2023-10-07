package downloader

import (
	"io/fs"
	"log"
	"os"
	"os/exec"
	"scout/Models"
	"sync"
	"time"

	"github.com/anacrolix/torrent"
)

var wg sync.WaitGroup

type Downloader struct {
	Client  *torrent.Client
	Torrent *torrent.Torrent
	File    Models.TorrentFile
}

func NewDownloader(directory string) *Downloader {
	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = directory
	client, _ := torrent.NewClient(cfg)

	return &Downloader{
		Client: client,
	}
}

func (d *Downloader) Init(requestedTitle string, scrapedFile Models.TorrentFile) {
	d.download(requestedTitle, scrapedFile)
}

func (d *Downloader) download(title string, file Models.TorrentFile) {
	t, err := d.Client.AddMagnet(file.Magnet)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Retrieving %s Torrent info ", title)
	<-t.GotInfo()
	log.Printf("%s information has been recieved", title)
	go monitor(title, t)
	start := time.Now()
	t.DownloadAll()
	if d.Client.WaitAll() {
		log.Printf("%s finished downloading , took %v", title, time.Since(start))
		defer t.Drop()
		moveRecentDownload(title)
		clearScreen()
	}
}

func monitor(title string, t *torrent.Torrent) {
	tick := time.NewTicker(5 * time.Second)
	for range tick.C {
		percentage := float64(t.BytesCompleted()) / float64(t.Info().TotalLength()) * 100
		clearScreen()
		log.Printf("%s Progress %.2f%%\n", title, percentage)
		if percentage == float64(100) {
			tick.Stop()
		}
	}
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func moveRecentDownload(title string) {
	workingDir := "/media/plex/downloads"
	downloads := getDownloads(workingDir)
	for _, download := range downloads {
		if download.IsDir() {
			found, err := os.Stat(workingDir + download.Name())
			if err != nil {
				log.Println(err)
			}
			sourcePath := workingDir + found.Name()
			destPath := "/media/plex/movies/" + title
			err = os.Rename(sourcePath, destPath)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func getDownloads(workingDir string) []fs.DirEntry {
	f, err := os.ReadDir(workingDir)
	if err != nil {
		log.Println(err)
	}
	return f
}
