package task

import models "scout/Models"

type Task struct {
	Name        string
	Sources     []string
	TorrentFile models.TorrentFile
}

func NewTask(name string, sources []string) *Task {
	return &Task{
		Name:        name,
		Sources:     sources,
		TorrentFile: models.TorrentFile{},
	}
}
