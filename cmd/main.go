package main

import (
	"context"
	"notion-helper/internal/app"
	"notion-helper/internal/helper"
	"notion-helper/internal/repository"
	"notion-helper/internal/service"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/jomei/notionapi"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func main() {
	helper.InitLog()

	e := echo.New()

	notion := notionapi.NewClient("your_integration_token")

	notionRepo := repository.NewNotionRepository(notion)

	notionService := service.NewNotionService(notionRepo)

	schedulerHandler := app.NewSchedulerHandler(notionService)

	s := gocron.NewScheduler(time.UTC)

	// TODO: change later
	s.Every(5).Seconds().Do(func() {
		schedulerHandler.FillEmptyTitleLinks(context.Background())
	})

	s.StartAsync()

	log.Error().Msg(e.Start(":4000").Error())
}
