package links

import (
	"context"

	"github.com/k1dan/crawler/internal/crawler"
	log "github.com/sirupsen/logrus"
)

type Crawler struct {
	crawler      *crawler.Crawler
	parsingRules *crawler.LinkParsingRules
}

// New func initialize new Crawler instance
func New(log *log.Logger, rules *crawler.LinkParsingRules) *Crawler {
	return &Crawler{
		crawler:      crawler.New(log),
		parsingRules: rules,
	}
}

// Run func runs links crawler
func (c *Crawler) Run(ctx context.Context, url string) (crawler.LinksParsingResult, error) {
	return c.crawler.GetItemLinks(ctx, url, c.parsingRules)
}
