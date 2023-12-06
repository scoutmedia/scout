package model

type TorrentFile struct {
	Name     string
	Size     float64
	Date     string
	Uploader string
	Magnet   string
}

type TorrentHTMLElement struct {
	Name     string
	Size     string
	Date     string
	Uploader string
	Magnet   string
	Hrefs    []string
}

type Media struct {
	Name         string    `json:"name,omitempty"`
	Title        string    `json:"title,omitempty"`
	Type         string    `json:"type,omitempty"`
	Overview     string    `json:"overview,omitempty"`
	OriginalName string    `json:"original_name,omitempty"`
	Adult        bool      `json:"adult,omitempty"`
	Keywords     []Keyword `json:"keywords,omitempty"`
}

type Keyword struct {
	Name string `json:"name,omitempty"`
	Id   int32  `json:"id,omitempty"`
}
