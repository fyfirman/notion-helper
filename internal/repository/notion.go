package repository

import (
	"context"
	"notion-helper/internal/datastruct"
	"notion-helper/internal/helper"
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

	query := notionapi.DatabaseQueryRequest{
		PageSize: 10,
	}

	db, err := r.notion.Database.Query(ctx, notionapi.DatabaseID(databaseId), &query)

	var notionLinks []datastruct.NotionLink

	for _, item := range db.Results {
		name := helper.RichTextsToString(item.Properties["Name"].(*notionapi.TitleProperty).Title)

		notionLinks = append(notionLinks, datastruct.NotionLink{
			Name: name,
			URL:  item.Properties["URL"].(*notionapi.URLProperty).URL,
			// TODO: add other props later
		})
	}

	if err != nil {
		return nil, err
	}

	return notionLinks, nil

}
