package scrapper

import (
	"container/list"
	"fmt"
	"log"
	"net/http"

	"github.com/mmcdole/gofeed"
	"github.com/tobg8/crypto-viz/common"
)

var maxCacheSize = 200 // Maximum number of article to keep in the cache
var processedLinks = list.New()

func isAlreadyProcessed(id string) bool {
	// Check if the link is in the cache
	for e := processedLinks.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == id {
			return true
		}
	}
	return false
}

func addToCache(link string) {
	// Check if the cache size exceeds the limit
	if processedLinks.Len() >= maxCacheSize {
		// Remove the oldest entry from the cache (the first one)
		oldest := processedLinks.Front()
		if oldest != nil {
			processedLinks.Remove(oldest)
		}
	}

	// Add the new link to the cache (at the end)
	processedLinks.PushBack(link)
}

func ScrapeRSSFeed(url string) ([]common.NewsEvent, error) {
	// Fetch the RSS feed
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS feed: %v", err)
	}
	defer response.Body.Close()

	// Parse the RSS feed
	fp := gofeed.NewParser()
	feed, err := fp.Parse(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed: %v", err)
	}

	var ne []common.NewsEvent
	for _, item := range feed.Items {
		// Check if the article's GUID has been processed before
		if isAlreadyProcessed(item.GUID) {
			log.Printf("already seen: %v", item.GUID)
			continue
		}

		image := ""
		if item.Image != nil {
			image = item.Image.URL
		}

		event := common.NewsEvent{
			ID:          item.GUID,
			Title:       item.Title,
			Link:        item.Link,
			RssURL:      url,
			ImageURL:    image,
			Author:      item.Authors[0].Name,
			PubDate:     *item.PublishedParsed,
			Categories:  item.Categories,
			Description: item.Description,
		}

		// Add the link to the cache
		addToCache(item.GUID)
		ne = append(ne, event)
	}

	return ne, nil
}
