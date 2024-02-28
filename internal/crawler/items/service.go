package items

import (
	"context"

	"github.com/k1dan/crawler/internal/crawler"
	log "github.com/sirupsen/logrus"
)

type Crawler struct {
	crawler      *crawler.Crawler
	parsingRules *crawler.ItemParsingRules
}

// New func initialize new Crawler instance
func New(log *log.Logger, rules *crawler.ItemParsingRules) *Crawler {
	return &Crawler{
		crawler:      crawler.New(log),
		parsingRules: rules,
	}
}

// Run func runs items crawler
func (c *Crawler) Run(ctx context.Context, url string) (crawler.Item, error) {
	return c.crawler.GetItemInfo(ctx, url, c.parsingRules)
}
