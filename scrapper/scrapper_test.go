package scrapper

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/tobg8/crypto-viz/common"
	"gotest.tools/assert"
)

type mockScraper struct {
	mock.Mock
}

// ScrapCurrencies is a mock implementation of the scraping function
func (m *mockScraper) ScrapCurrencies(url string) ([]common.Currency, error) {
	args := m.Called(url)
	return args.Get(0).([]common.Currency), args.Error(1)
}

// ScrapCurrencies is a mock implementation of the scraping function
func (m *mockScraper) ScrapNewCurrencies(url string) ([]common.CurrencyEvent, error) {
	args := m.Called(url)
	return args.Get(0).([]common.CurrencyEvent), args.Error(1)
}

func Test_ScrapCurrencies(t *testing.T) {
	type test struct {
		wantError bool
		url       string
		want      []common.Currency
		mockCall  func(m *mockScraper)
	}

	tests := map[string]test{
		"unexpected length": {
			url:  "https://www.coingecko.com/fr",
			want: []common.Currency{},
			mockCall: func(m *mockScraper) {
				m.On("ScrapCurrencies", "https://www.coingecko.com/fr").
					Return(
						[]common.Currency{}, errors.New("some error"),
					)
			},
			wantError: true,
		},
		"bad url": {
			url: "hs://www.coingecko.com/fr",
			mockCall: func(m *mockScraper) {
				m.On("ScrapCurrencies", "hs://www.coingecko.com/fr").
					Return(
						[]common.Currency{}, errors.New("some error"),
					)
			},
			wantError: true,
		},
		"nominal": {
			url: "https://www.coingecko.com/fr",
			want: []common.Currency{
				{
					ID:   "1",
					Name: "MockCurrency1",
				},
				{
					ID:   "2",
					Name: "MockCurrency2",
				},
			},
			mockCall: func(m *mockScraper) {
				m.On("ScrapCurrencies", "https://www.coingecko.com/fr").
					Return(
						[]common.Currency{
							{
								ID:   "1",
								Name: "MockCurrency1",
							},
							{
								ID:   "2",
								Name: "MockCurrency2",
							},
						}, nil)
			},
			wantError: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockScraper := new(mockScraper)

			if tt.mockCall != nil {
				tt.mockCall(mockScraper)
			}

			currencies, err := mockScraper.ScrapCurrencies(tt.url)
			mockScraper.AssertExpectations(t)
			if tt.wantError {
				assert.Error(t, err, "some error")
				return
			}

			assert.NilError(t, err)
			assert.DeepEqual(t, currencies, tt.want)
		})
	}
}

func Test_ScrapNewCurrencies(t *testing.T) {
	type test struct {
		wantError bool
		url       string
		want      []common.CurrencyEvent
		mockCall  func(m *mockScraper)
	}

	tests := map[string]test{
		"unexpected length": {
			wantError: true,
			url:       "https://www.coingecko.com/fr/new-cryptocurrencies",
			mockCall: func(m *mockScraper) {
				m.On("ScrapNewCurrencies", "https://www.coingecko.com/fr/new-cryptocurrencies").
					Return(
						[]common.CurrencyEvent{}, errors.New("some error"),
					)
			},
		},
		"bad url": {
			wantError: true,
			url:       "hps://www.coingecko.com/fr/new-cryptocurrencies",
			mockCall: func(m *mockScraper) {
				m.On("ScrapNewCurrencies", "hps://www.coingecko.com/fr/new-cryptocurrencies").
					Return(
						[]common.CurrencyEvent{}, errors.New("some error"),
					)
			},
		},
		"nominal": {
			wantError: false,
			url:       "https://www.coingecko.com/fr/new-cryptocurrencies",
			want: []common.CurrencyEvent{
				{
					Name: "MockCurrency1",
				},
				{
					Name: "MockCurrency2",
				},
			},
			mockCall: func(m *mockScraper) {
				m.On("ScrapNewCurrencies", "https://www.coingecko.com/fr/new-cryptocurrencies").
					Return(
						[]common.CurrencyEvent{
							{
								Name: "MockCurrency1",
							},
							{
								Name: "MockCurrency2",
							},
						}, nil)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockScraper := new(mockScraper)

			if tt.mockCall != nil {
				tt.mockCall(mockScraper)
			}

			currencies, err := mockScraper.ScrapNewCurrencies(tt.url)
			mockScraper.AssertExpectations(t)
			if tt.wantError {
				assert.Error(t, err, "some error")
				return
			}

			assert.NilError(t, err)
			assert.DeepEqual(t, currencies, tt.want)
		})
	}
}
