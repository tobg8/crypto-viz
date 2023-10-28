package usecase

import (
	"fmt"
	"log"
	"strings"

	"github.com/tobg8/crypto-viz/common"
	"github.com/tobg8/crypto-viz/internal/repository"
)

var lastURL string

// HandleNews will fetch news from API transform and send them through kafka
func HandleNews(k *repository.KafkaClient) error {
	articles := repository.FetchArticles()
	if len(*articles) == 0 {
		return fmt.Errorf("no news to send")
	}

	var events []common.ArticleEvent
	for _, v := range *articles {
		event := transformArticleToEvent(v)
		events = append(events, event)
	}

	newEvents := filterSentEvents(events, &lastURL)
	if len(newEvents) == 0 {
		return nil
	}
	lastURL = newEvents[0].URL

	// Push News
	err := k.PushArticles(newEvents)
	if err != nil {
		return fmt.Errorf("could not send news: %w", err)
	}

	return nil
}

// filterSentEvents returns the events not processed by the producer
func filterSentEvents(ae []common.ArticleEvent, url *string) []common.ArticleEvent {
	var newEvents []common.ArticleEvent
	for _, event := range ae {
		if event.URL != *url {
			newEvents = append(newEvents, event)
		} else {
			log.Print("There is no articles to send yet")
			break
		}
	}

	return newEvents
}

// transformArticleToEvent transform article fields
func transformArticleToEvent(a common.ArticleAPI) common.ArticleEvent {
	currencies := TransformCurrency(a.Currencies)
	return common.ArticleEvent{
		Kind:        a.Kind,
		Source:      a.SourceAPI.Domain,
		Title:       a.Title,
		URL:         a.SourceAPI.Domain + "/" + a.Slug,
		Currencies:  currencies,
		PublishedAt: a.PublishedAt,
		CreatedAt:   a.CreatedAt,
	}
}

// TransformCurrency transform currency fields
func TransformCurrency(c []common.CurrencyAPI) []common.CurrencyEvent {
	var currencies []common.CurrencyEvent
	for _, v := range c {
		currencies = append(currencies, common.CurrencyEvent{
			Code:  strings.ToLower(v.Code),
			Title: v.Slug,
		})
	}

	return currencies
}
