package model

type TorrentFile struct {
	Name     string
	Size     float64
	Date     string
	Uploader string
	Magnet   string
}

type Request struct {
	Data string `json:"data,omitempty"`
}
