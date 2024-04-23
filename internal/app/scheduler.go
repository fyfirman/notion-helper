package app

import (
	"context"
	"notion-helper/internal/service"

	"github.com/rs/zerolog/log"
)

type SchedulerHandlerInterface interface {
	FillEmptyTitleLinks(ctx context.Context) error
}

type SchedulerHandler struct {
	notionService service.NotionServiceInterface

	SchedulerHandlerInterface
}

func NewSchedulerHandler(notionService service.NotionServiceInterface) *SchedulerHandler {
	return &SchedulerHandler{
		notionService: notionService,
	}
}

func (h *SchedulerHandler) FillEmptyTitleLinks(ctx context.Context) {
	log.Info().Msg("Check empty title links scheduler")
	err := h.notionService.FillEmptyTitleLinks(ctx)

	if err != nil {
		log.Error().Msg(err.Error())
	}
}
