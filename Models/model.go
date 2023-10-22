package model

type TorrentFile struct {
	Name     string
	Size     float64
	Date     string
	Uploader string
	Magnet   string
}

type Media struct {
	Data struct {
		Adult        bool   `json:"adult,omitempty"`
		Overview     string `json:"overview,omitempty"`
		OriginalName string `json:"original_name,omitempty"`
		Name         string `json:"name,omitempty"`
		Type         string `json:"type,omitempty"`
		Keywords     []struct {
			Name string `json:"name,omitempty"`
			Id   int32  `json:"id,omitempty"`
		} `json:"keywords,omitempty"`
	} `json:"data,omitempty"`
}
