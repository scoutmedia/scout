package downloader

import (
	"io/fs"
	"log"
	"os"
	"os/exec"
	model "scout/Models"
	"sync"
	"time"

	"github.com/anacrolix/torrent"
)

var wg sync.WaitGroup

type Downloader struct {
	Client  *torrent.Client
	Torrent *torrent.Torrent
}

func NewDownloader(dataDir string) *Downloader {
	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = dataDir
	client, _ := torrent.NewClient(cfg)
	return &Downloader{
		Client: client,
	}
}

func (d *Downloader) Start(title string, torrentFile model.TorrentFile) {
	t, err := d.Client.AddMagnet(torrentFile.Magnet)
	if err != nil {
		log.Println("Error occured adding torrent magnet", err)
	}
	log.Printf("Retrieving %s info", torrentFile.Name)
	<-t.GotInfo()
	log.Printf("%s retrived info", torrentFile.Name)
	log.Printf("%s download will begin shortly ..", torrentFile.Name)
	start := time.Now()
	t.DownloadAll()
	if d.Client.WaitAll() {
		log.Printf("%s finished downloading , took %v", torrentFile.Name, time.Since(start))
		defer t.Drop()
		clearScreen()
		moveRecentDownload(title)
	}
}

func (d *Downloader) Monitor(client *torrent.Client) {
	tick := time.NewTicker(5 * time.Second)
	for range tick.C {
		if len(client.Torrents()) != 0 {
			log.Println("Active Downloads")
			go func() {
				time.Sleep(9 * time.Second)
				clearScreen()
			}()
			for _, activeTorrent := range client.Torrents() {
				go func(activeTorrent *torrent.Torrent) {
					info := activeTorrent.Info()
					if info != nil {
						torrent, _ := d.Client.Torrent(activeTorrent.InfoHash())
						percentage := float64(torrent.BytesCompleted()) / float64(torrent.Info().TotalLength()) * 100
						log.Printf("%s Progress %.2f%%\n", torrent.Info().Name, percentage)
					}
				}(activeTorrent)
			}
		}
	}
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func moveRecentDownload(title string) {
	workingDir := "/media/plex/downloads/"
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
