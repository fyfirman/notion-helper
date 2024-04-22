package service

import (
	"context"
	"notion-helper/internal/repository"
)

type NotionServiceInterface interface {
}

type NotionService struct {
	notionRepo repository.NotionRepositoryInterface

	NotionServiceInterface
}

func NewNotionService(
	notionRepo repository.NotionRepositoryInterface,
) *NotionService {
	return &NotionService{
		notionRepo: notionRepo,
	}
}

func (s *NotionService) FillEmptyTitleLinks(ctx context.Context) {
	// TODO: later
}
