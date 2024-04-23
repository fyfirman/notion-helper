package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"notion-helper/internal/datastruct"
	"notion-helper/internal/repository"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	html "github.com/levigross/exp-html"
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

	log.Info().Msgf("Got %d links from API", len(links))

	emptyTitleLinks := make([]datastruct.NotionLink, 0)
	for _, link := range links {
		if link.Name == "" && link.URL != "" {
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

			err := s.FillEmptyTitle(ctx, l.URL)

			if err != nil {
				log.Error().Str("link", l.URL).Err(err).Msg("Error while filling empty title links")
			}
			resultChan <- err
		}(link)
	}

	wg.Wait()

	close(resultChan)

	return nil
}

func (s *NotionService) FillEmptyTitle(ctx context.Context, link string) error {
	title, err := s.getSeoTitle(ctx, link)

	if err != nil {
		log.Error().Err(err).Msg("Error while filling empty title")
		return err
	}

	log.Info().Str("link", link).Str("title", title).Msgf("Got SEO title")

	// TODO: P0 - update to database

	return nil
}

func (s *NotionService) getSeoTitle(ctx context.Context, link string) (string, error) {
	unsupportedDomain := []string{
		"x.com",
		"twitter.com",
	}

	isUnsupported := false

	for _, domain := range unsupportedDomain {
		if strings.Contains(link, domain) {
			isUnsupported = true
		}
	}

	if isUnsupported {
		return s.getSeoTitleWithChrome(ctx, link)
	}

	resp, err := http.Get(link)
	if err != nil {
		log.Error().Err(err).Msg("Error while get seo title")
		return "", err
	}
	defer resp.Body.Close()

	return parseHTMLForTitle(resp.Body)
}

func parseHTMLForTitle(r io.Reader) (string, error) {
	z := html.NewTokenizer(r)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			log.Error().Err(z.Err()).Msg("Error token while parsing")
			return "", z.Err()
		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()
			if t.Data == "title" {
				tt = z.Next()
				if tt == html.TextToken {
					return z.Token().Data, nil
				}
			}
		}
	}

	return "", nil // ASK: Why this is unreachable??
}

func (s *NotionService) getSeoTitleWithChrome(ctx context.Context, link string) (string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, 10*time.Second)
	defer cancelTimeout()

	var title, htmlContent string
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(link), // Navigate to the page
		// chromedp.WaitVisible(`title`, chromedp.ByQuery), // Wait until the title element is visible
		chromedp.OuterHTML("html", &htmlContent, chromedp.ByQuery),
		// chromedp.ActionFunc(func(ctx context.Context) error {
		// 	for {
		// 		err := chromedp.Title(&title).Do(ctx)
		// 		if err != nil {
		// 			return err
		// 		}
		// 		if strings.TrimSpace(title) != "" {
		// 			break
		// 		}
		// 		time.Sleep(500 * time.Millisecond)
		// 	}
		// 	return nil
		// }),
	)

	fmt.Println(htmlContent)

	if err != nil {
		log.Error().Err(err).Str("HTML", htmlContent).Msg("Failed to retrieve title")
		return "", err
	}

	return title, nil

}

// QnA
// Tanya Idzhar
//  1. kenapa function yg gak di implement gak ada error
// 2. Apakah harus Msg("") ?
