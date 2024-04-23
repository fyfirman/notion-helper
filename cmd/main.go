package main

import (
	"notion-helper/internal/app"
	"notion-helper/internal/helper"
	"notion-helper/internal/repository"
	"notion-helper/internal/service"

	"github.com/jomei/notionapi"
)

func main() {
	helper.InitLog()

	notion := notionapi.NewClient("your_integration_token")

	notionRepo := repository.NewNotionRepository(notion)

	notionService := service.NewNotionService(notionRepo)

	scheduler := app.NewScheduler(notionService)

}
