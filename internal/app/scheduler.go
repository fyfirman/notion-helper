package app

import "notion-helper/internal/service"

type SchedulerInterface interface {
}

type Scheduler struct {
	notionService service.NotionServiceInterface

	SchedulerInterface
}

func NewScheduler(notionService service.NotionServiceInterface) *Scheduler {
	return &Scheduler{
		notionService: notionService,
	}
}
