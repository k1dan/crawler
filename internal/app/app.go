package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/k1dan/crawler/internal/config"
	"github.com/k1dan/crawler/internal/crawler"
	"github.com/k1dan/crawler/internal/crawler/items"
	"github.com/k1dan/crawler/internal/crawler/links"
	"github.com/k1dan/crawler/internal/storage"
	"github.com/k1dan/crawler/internal/storage/jsonstorage"
	"github.com/k1dan/crawler/pkg/workerpool"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type App struct {
	logger              *log.Logger
	linksParsingRules   *crawler.LinkParsingRules
	itemParsingRules    *crawler.ItemParsingRules
	ItemsParsingWorkers workerpool.WorkerPool
	SavePath            string
}

func New(logger *log.Logger, config *config.Config) *App {
	cdnToParse := map[string]struct{}{}
	for _, c := range config.ConditionsToParse {
		cdnToParse[c] = struct{}{}
	}
	// TODO: refactor using crawler constructor
	// TODO: use json to define parsing rules
	itemParsingRules := &crawler.ItemParsingRules{
		Title: crawler.Element{
			Tag:   "div",
			Class: "x-item-title",
		},
		Condition: crawler.Element{
			Tag:   "div",
			Class: "x-item-condition-text",
		},
		Price: crawler.Element{
			Tag:   "div",
			Class: "x-price-primary",
		},
		ID: crawler.Element{
			Tag:   "div",
			Class: "ux-layout-section__textual-display--itemId",
		},
		ConditionsToParse: cdnToParse,
	}
	linkParsingRules := &crawler.LinkParsingRules{
		ItemURL: crawler.Element{
			Tag:           "a",
			Class:         "s-item__link",
			AttributeName: "href",
		},
		NextPageLink: crawler.Element{
			Tag:           "a",
			Class:         "pagination__next",
			AttributeName: "href",
		},
	}
	return &App{
		logger:              logger,
		linksParsingRules:   linkParsingRules,
		itemParsingRules:    itemParsingRules,
		ItemsParsingWorkers: workerpool.New(config.ParserWorkerAmount),
		SavePath:            config.FileSavePath,
	}
}

func (app *App) Run(ctx context.Context, startURL string) error {
	app.logger.Infof("start crawling")
	var errChan chan error
	var wg sync.WaitGroup

	go app.ItemsParsingWorkers.Run(ctx)

	jsonStorage := jsonstorage.New(app.SavePath)
	jsonStorage.Listen(ctx, &wg)

	go app.startCrawl(ctx, startURL, errChan)

	for {
		select {
		case r, ok := <-app.ItemsParsingWorkers.Results():
			if !ok {
				app.logger.Infof("finishing crawling")
				jsonStorage.Close()
				wg.Wait()
				return nil
			}
			log.Infof("results from job %v has been received, value: %v", r.Description.ID, r.Value)
			if r.Err != nil {
				log.Errorf("error from job %v has been received", r.Err)
			}
			item, ok := r.Value.(crawler.Item)
			if !ok {
				log.Errorf("result from job %v can not be casted", r.Description.ID)
				continue
			}

			jsonStorage.Save(storage.File{
				FileName: item.ID,
				Item: storage.Item{
					Title:      item.Title,
					Condition:  item.Condition,
					Price:      item.Price,
					ProductURL: item.ProductURL,
				},
			})
		case err := <-errChan:
			return fmt.Errorf("run error: %w", err)
		default:
		}
	}
}

func (app *App) startCrawl(ctx context.Context, startURL string, errChan chan error) {
	linksCrawler := links.New(app.logger, app.linksParsingRules)
	var linksCrawlingResult crawler.LinksParsingResult
	var err error
	// TODO: refactor NexPageLink = "1"
	linksCrawlingResult.NexPageLink = "1"
	for linksCrawlingResult.NexPageLink != "" {
		linksCrawlingResult, err = linksCrawler.Run(ctx, startURL)
		if err != nil {
			app.logger.Errorf("failed to get links: %v", err)
			errChan <- err
		}
		for _, l := range linksCrawlingResult.Links {
			itemsCrawler := items.New(app.logger, app.itemParsingRules)

			execFunc := func(ctx context.Context, arg any) (any, error) {
				return itemsCrawler.Run(ctx, fmt.Sprintf("%v", arg))
			}
			job := workerpool.Job{
				Description: workerpool.JobDescription{
					ID:   workerpool.JobID(uuid.New().String()),
					Type: workerpool.JobType("item crawler"),
				},
				ExecFunc: execFunc,
				Args:     l,
			}
			app.ItemsParsingWorkers.AddJob(job)

			// needed to avoid blocking from ebay
			time.Sleep(time.Second * 1)
		}
	}
	app.ItemsParsingWorkers.CloseJobsChan()
}
