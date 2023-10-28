package usecase

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tobg8/crypto-viz/common"
)

func TestFiltrerSentEvent(t *testing.T) {
	t.Run("when slug is not present", func(t *testing.T) {
		in := []common.ArticleEvent{
			{
				Title: "nICE",
				URL:   "r",
			},
			{
				Title: "Hello",
				URL:   "url-super",
			},
			{
				Title: "Yaho",
				URL:   "url-second",
			},
		}

		slug := ""

		expect := []common.ArticleEvent{
			{
				Title: "nICE",
				URL:   "r",
			},
			{
				Title: "Hello",
				URL:   "url-super",
			},
			{
				Title: "Yaho",
				URL:   "url-second",
			},
		}

		out := filterSentEvents(in, &slug)
		assert.Equal(t, expect, out)
	})

	t.Run("when slug is present", func(t *testing.T) {
		in := []common.ArticleEvent{
			{
				Title: "nICE",
				URL:   "r",
			},
			{
				Title: "Hello",
				URL:   "url-super",
			},
			{
				Title: "Yaho",
				URL:   "url-second",
			},
		}

		slug := "url-super"

		expect := []common.ArticleEvent{
			{
				Title: "nICE",
				URL:   "r",
			},
		}

		out := filterSentEvents(in, &slug)
		assert.Equal(t, expect, out)
	})
}

func TestTransformArticleToEvent(t *testing.T) {
	in := common.ArticleAPI{
		Kind: "media",
		SourceAPI: common.SourceAPI{
			Title:  "source title",
			Region: "fr",
			Domain: "kia.com",
		},
		Title: "title",
		URL:   "http://crypto",
		Slug:  "slugify-slug",
		Currencies: []common.CurrencyAPI{
			{
				Code:  "BNB",
				Title: "BNB",
				Slug:  "binancecoin",
				URL:   "https://cryptopanic.com/news/binancecoin/",
			},
			{
				Code:  "BTC",
				Title: "Bitcoin",
				Slug:  "bitcoin",
				URL:   "https://cryptopanic.com/news/bitcoin/",
			},
		},
		PublishedAt: time.Time{},
		CreatedAt:   time.Time{},
	}

	expect := common.ArticleEvent{
		Kind:   "media",
		Source: "kia.com",
		Title:  "title",
		URL:    "kia.com/slugify-slug",
		Currencies: []common.CurrencyEvent{
			{
				Code:  "bnb",
				Title: "binancecoin",
			},
			{
				Code:  "btc",
				Title: "bitcoin",
			},
		},
		PublishedAt: time.Time{},
		CreatedAt:   time.Time{},
	}

	out := transformArticleToEvent(in)
	assert.Equal(t, expect, out)
}

func TestTransformCurrency(t *testing.T) {
	in := []common.CurrencyAPI{
		{
			Code:  "BNB",
			Title: "BNB",
			Slug:  "binancecoin",
			URL:   "https://cryptopanic.com/news/binancecoin/",
		},
		{
			Code:  "BTC",
			Title: "Bitcoin",
			Slug:  "bitcoin",
			URL:   "https://cryptopanic.com/news/bitcoin/",
		},
	}

	expect := []common.CurrencyEvent{
		{
			Code:  "bnb",
			Title: "binancecoin",
		},
		{
			Code:  "btc",
			Title: "bitcoin",
		},
	}

	out := TransformCurrency(in)
	assert.True(t, reflect.DeepEqual(expect, out))
}
