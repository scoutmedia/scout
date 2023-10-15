package downloader

import (
	"fmt"
	"log"
	"os"
	logger "scout/Logger"
	model "scout/Models"
	"time"

	"github.com/anacrolix/torrent"
)

type Downloader struct {
	Client  *torrent.Client
	Torrent *torrent.Torrent
	logger  *logger.Logger
}

func NewDownloader(dataDir string, logger *logger.Logger) *Downloader {
	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = dataDir
	client, _ := torrent.NewClient(cfg)
	return &Downloader{
		Client: client,
		logger: logger,
	}
}

func (d *Downloader) Start(title string, torrentFile model.TorrentFile) {
	t, err := d.Client.AddMagnet(torrentFile.Magnet)
	if err != nil {
		d.logger.Error("Download", fmt.Sprint("Error occured adding torrent magnet", err))
	}
	d.logger.Info("Download", fmt.Sprintf("Retrieving %s info", torrentFile.Name))
	<-t.GotInfo()
	d.logger.Info("Download", fmt.Sprintf("%s retrived info", torrentFile.Name))
	d.logger.Info("Download Start", fmt.Sprintf("%s download will begin shortly...", torrentFile.Name))
	go d.status(title, t)
	t.DownloadAll()
}

func (d *Downloader) status(title string, t *torrent.Torrent) {
	start := time.Now()
	tick := time.NewTicker(5 * time.Second)
	for range tick.C {
		activeTorrent, _ := d.Client.Torrent(t.InfoHash())
		percentage := float64(activeTorrent.BytesCompleted()) / float64(activeTorrent.Info().TotalLength()) * 100
		log.Printf("%s Progress %.2f%%\n", activeTorrent.Name(), percentage)

		if percentage == 100 {
			log.Printf("%s download took %v", title, time.Since(start))
			d.logger.Info("Download Complete", fmt.Sprintf("%s download took %v", title, time.Since(start)))
			d.moveRecentDownload(title, t.Info().Name)
			tick.Stop()
			d.drop(t)
		}
	}
}

func (d *Downloader) drop(t *torrent.Torrent) {
	t.Drop()
}

func (d *Downloader) moveRecentDownload(title string, folderName string) {
	os.Rename("/media/plex/downloads/"+folderName, "/media/plex/movies/"+title)
}
