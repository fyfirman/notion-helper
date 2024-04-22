package repository

import (
	"context"

	"github.com/jomei/notionapi"
)

type NotionRepositoryInterface interface {
}

type NotionRepository struct {
	notion notionapi.Client

	NotionRepositoryInterface
}

func NewNotionRepository(notion *notionapi.Client) *NotionRepository {
	return &NotionRepository{
		notion: *notion,
	}
}

func (r NotionRepository) GetAllLinks(ctx context.Context) {
	// TODO: later
}
