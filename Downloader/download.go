package downloader

import (
	"log"
	"os"
	model "scout/Models"
	"time"

	"github.com/anacrolix/torrent"
)

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
	log.Println(torrentFile.Name, t.Info().Name)
	go d.status(title, t)
	t.DownloadAll()

	// if d.Client.WaitAll() {
	// 	log.Printf("%s finished downloading , took %v", torrentFile.Name, time.Since(start))
	// 	defer t.Drop()
	// 	d.moveRecentDownload(title, t.Info().Name)
	// }
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
