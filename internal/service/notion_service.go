package service

import (
	"context"
	"notion-helper/internal/datastruct"
	"notion-helper/internal/repository"
	"sync"

	"github.com/rs/zerolog/log"
)

type NotionServiceInterface interface {
	FillEmptyTitleLinks(ctx context.Context) error
	FillEmptyTitle(ctx context.Context, link string) error
	getSeoTitle(ctx context.Context, link string) (string, error)
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

func (s *NotionService) FillEmptyTitleLinks(ctx context.Context) error {
	links, err := s.notionRepo.GetAllLinks(ctx)

	if err != nil {
		return err
	}

	if len(links) == 0 {
		log.Info().Msg("Got 0 links from API")
		return nil
	}

	log.Info().Msgf("Got %s links from API", len(links))

	emptyTitleLinks := make([]datastruct.NotionLink, 0)
	for _, link := range links {
		if link.Title != "" {
			emptyTitleLinks = append(emptyTitleLinks, link)
		}
	}

	resultChan := make(chan error, len(emptyTitleLinks))

	log.Info().Msgf("%d links has empty titles", len(emptyTitleLinks))

	var wg sync.WaitGroup

	for _, link := range emptyTitleLinks {
		wg.Add(1)
		go func(l datastruct.NotionLink) {
			defer wg.Done()

			err := s.FillEmptyTitle(ctx, l.Title)

			log.Error().Str("link", l.Title).Msg(err.Error())
			resultChan <- err
		}(link)
	}

	wg.Wait()

	close(resultChan)

	return nil
}
