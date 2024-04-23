package main

import (
	"context"
	"notion-helper/internal/app"
	"notion-helper/internal/helper"
	"notion-helper/internal/repository"
	"notion-helper/internal/service"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/jomei/notionapi"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func main() {
	helper.InitLog()

	err := godotenv.Load(".env")

	if err != nil {
		panic(err.Error())
	}

	e := echo.New()

	notionToken := os.Getenv("NOTION_TOKEN")

	notion := notionapi.NewClient(notionapi.Token(notionToken))

	notionRepo := repository.NewNotionRepository(notion)

	notionService := service.NewNotionService(notionRepo)

	schedulerHandler := app.NewSchedulerHandler(notionService)

	s := gocron.NewScheduler(time.UTC)

	s.Cron(os.Getenv("CRON_FILLING_EMPTY_LINKS")).Do(func() {
		schedulerHandler.FillEmptyTitleLinks(context.Background())
	})

	s.StartAsync()

	log.Error().Msg(e.Start(":4000").Error())
}
