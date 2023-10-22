package task

import models "scout/Models"

type Task struct {
	Media       models.Media
	Sources     []string
	TorrentFile models.TorrentFile
}

func NewTask(media models.Media, sources []string) *Task {
	return &Task{
		Media:       media,
		Sources:     sources,
		TorrentFile: models.TorrentFile{},
	}
}
