package crawler

import (
	"fmt"
)

const ebayPrefix = "https://www.ebay.com/"

type (
	// Element TODO: make proper constructor
	Element struct {
		Tag           string
		Class         string
		ID            string
		AttributeName string
	}
	LinkParsingRules struct {
		ItemURL      Element
		NextPageLink Element
	}
	LinksParsingResult struct {
		Links       []string
		NexPageLink string
	}
	ItemParsingRules struct {
		Title             Element
		Condition         Element
		Price             Element
		ID                Element
		ConditionsToParse map[string]struct{}
	}
	Item struct {
		Title      string
		Condition  string
		Price      string
		ProductURL string
		ID         string
	}
)

// BuildHTMLSelector func builds valid html selector from tag, class and id
func (e Element) BuildHTMLSelector() string {
	var selector string
	if e.Tag != "" {
		selector = e.Tag
	}
	if e.Class != "" {
		selector = fmt.Sprintf("%s.%s", selector, e.Class)
	}
	if e.ID != "" {
		selector = fmt.Sprintf("%s#%s", selector, e.ID)
	}
	return selector
}

// GetID func is ID getter
func (i Item) GetID() string {
	return i.ID
}
