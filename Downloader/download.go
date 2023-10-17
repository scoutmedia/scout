package downloader

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
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
		return
	}
	d.logger.Info("Download", fmt.Sprintf("Retrieving %s info", torrentFile.Name))
	<-t.GotInfo()
	d.logger.Info("Download", fmt.Sprintf("%s retrived info", torrentFile.Name))
	d.logger.Info("Download", fmt.Sprintf("%s download will begin shortly...", torrentFile.Name))
	go d.status(title, t)
	t.DownloadAll()
}

func (d *Downloader) status(title string, t *torrent.Torrent) {
	start := time.Now()
	tick := time.NewTicker(10 * time.Second)
	for range tick.C {
		activeTorrent, _ := d.Client.Torrent(t.InfoHash())
		percentage := float64(activeTorrent.BytesCompleted()) / float64(activeTorrent.Info().TotalLength()) * 100
		log.Printf("%s Progress %.2f%%\n", activeTorrent.Name(), percentage)

		if percentage == 100 {
			tick.Stop()
			d.logger.Info("Download Complete", fmt.Sprintf("%s download took %v", title, time.Since(start)))
			d.moveFile(t.Info().Name, title)
		}
	}
}

func (d *Downloader) moveFile(source string, fileName string) {
	fs, err := os.ReadDir("/media/plex/downloads/" + source)
	if err != nil {
		d.logger.Error("moveFile", fmt.Sprintf("Error occured accessing %s directory", source))
	}

	for _, f := range fs {
		if path.Ext(f.Name()) == ".mkv" || path.Ext(f.Name()) == ".mp4" {
			fileSrc, err := os.Open(fmt.Sprintf("/media/plex/downloads/%s/%s", source, f.Name()))
			if err != nil {
				d.logger.Error("moveFile", fmt.Sprintf("Error opening media file %s", err))
			}
			err = os.Mkdir(fmt.Sprintf("/media/plex/movies/%s", fileName), 0775)
			if err != nil {
				d.logger.Error("moveFile", fmt.Sprintf("Error creating %s directory: %v", fileName, err))
			}
			dst, err := os.Create(fmt.Sprintf("/media/plex/movies/%s/%s%s", fileName, source, path.Ext(fileSrc.Name())))
			if err != nil {
				d.logger.Error("moveFile", fmt.Sprintf("Error creating %s movie file: %v", fileName, err))
			}
			_, err = io.Copy(dst, fileSrc)
			if err != nil {
				d.logger.Error("moveFile", fmt.Sprintf("Error occured copying %s data to destination : %s", fileSrc.Name(), dst.Name()))
			}
		}
	}
}
