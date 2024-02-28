package crawler

import (
	"context"
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

const bodyTag = "body"

type Crawler struct {
	logger    *log.Logger
	collector *colly.Collector
}

// New initializes new Crawler instance
func New(logger *log.Logger) *Crawler {
	return &Crawler{
		logger:    logger,
		collector: colly.NewCollector(),
	}
}

// GetItemLinks parses item links using LinkParsingRules
func (c *Crawler) GetItemLinks(
	ctx context.Context,
	url string,
	parsingRules *LinkParsingRules,
) (LinksParsingResult, error) {
	var result LinksParsingResult

	c.collector.OnHTML(bodyTag, func(body *colly.HTMLElement) {
		result.Links = body.ChildAttrs(
			parsingRules.ItemURL.BuildHTMLSelector(),
			parsingRules.ItemURL.AttributeName,
		)
		result.NexPageLink = body.ChildAttr(
			parsingRules.NextPageLink.BuildHTMLSelector(),
			parsingRules.NextPageLink.AttributeName,
		)
	})

	err := c.start(ctx, url)
	if err != nil {
		return result, fmt.Errorf("GetItemsLink: %w", err)
	}

	return result, nil
}

// GetItemInfo parses item info using ItemParsingRules
func (c *Crawler) GetItemInfo(
	ctx context.Context,
	url string,
	parsingRules *ItemParsingRules,
) (Item, error) {
	var result Item
	result.ProductURL = url

	c.collector.OnHTML(bodyTag, func(body *colly.HTMLElement) {
		result.Title = body.ChildText(parsingRules.Title.BuildHTMLSelector())
		result.Price = body.ChildText(parsingRules.Price.BuildHTMLSelector())
		condition := ""
		// TODO: refactor
		body.ForEach(parsingRules.Condition.BuildHTMLSelector(), func(_ int, e *colly.HTMLElement) {
			condition = e.ChildText("span.clipped")
		})
		body.ForEach(parsingRules.ID.BuildHTMLSelector(), func(_ int, e *colly.HTMLElement) {
			result.ID = e.ChildText("span.ux-textspans--BOLD")
		})
		if result.ID == "" {
			replace := strings.ReplaceAll(url, ebayPrefix, "")
			split := strings.SplitAfterN(replace, "/", 2)
			id := split[0]
			if len(split) > 1 {
				id = strings.SplitAfterN(split[1], "?", 2)[0]
			}
			result.ID = strings.ReplaceAll(id, "?", "")
		}
		if len(parsingRules.ConditionsToParse) != 0 {
			_, ok := parsingRules.ConditionsToParse[condition]
			if !ok {
				return
			}
			result.Condition = condition
		}
		result.Condition = condition
	})

	err := c.start(ctx, url)
	if err != nil {
		log.Errorf("getIemINFO err: %w", err)
		return result, fmt.Errorf("GetItemInfo: %w", err)
	}
	if result.Condition == "" {
		return Item{}, nil
	}
	log.Infof("result: %v", result)
	return result, nil
}

// start function starts crawling the page
func (c *Crawler) start(ctx context.Context, url string) error {
	c.collector.OnRequest(func(r *colly.Request) {
		c.logger.WithContext(ctx).
			Infof("Visiting: %v", r.URL.String())
	})
	c.collector.OnError(func(r *colly.Response, err error) {
		c.logger.WithContext(ctx).
			Errorf("Error visiting page: %v; error description %v", r.Request.URL.String(), err)
	})

	err := c.collector.Visit(url)
	if err != nil {
		return fmt.Errorf("error visiting page: %w", err)
	}
	return nil
}
