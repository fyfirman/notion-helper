package repository

import (
	"context"
	"notion-helper/internal/datastruct"
	"os"

	"github.com/jomei/notionapi"
)

type NotionRepositoryInterface interface {
	GetAllLinks(ctx context.Context) ([]datastruct.NotionLink, error)
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

func (r NotionRepository) GetAllLinks(ctx context.Context) ([]datastruct.NotionLink, error) {
	databaseId := os.Getenv("NOTION_DATABASE_ID")

	_, err := r.notion.Database.Get(ctx, notionapi.DatabaseID(databaseId))

	if err != nil {
		return nil, err
	}

	return nil, nil

}
