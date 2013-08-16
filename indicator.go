package wbdata

import (
	"fmt"
)

// Indicator provides access to the Indicator related queries
// in the World Bank Open Data API.
//
// World Bank Open Data API docs: http://data.worldbank.org/node/203
type IndicatorService struct {
	client *Client
}

type Indicator struct {
	Id         string
	Name       string
	Source     *Source
	SourceNote string
}

func (i *IndicatorService) ListIndicators() ([]Indicator, error) {
	indicator := []Indicator{}

	req, err := i.client.NewRequest("GET", "indicators", nil)
	if err != nil {
		return nil, err
	}

	_, err = i.client.Do(req, &[]interface{}{&indicator})

	return indicator, err

}

func (i *IndicatorService) GetIndicator(id string) ([]Indicator, error) {
	indicator := []Indicator{}
	urlStr := fmt.Sprintf("indicators/%v", id)
	req, err := i.client.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	_, err = i.client.Do(req, &[]interface{}{&indicator})

	return indicator, err
}
