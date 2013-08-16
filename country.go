package wbdata

import (
	"fmt"
)

// CountryService provides access to the countries related queries
// in the World Bank Open Data API.
//
// World Bank Open Data API docs: http://data.worldbank.org/node/18
type CountryService struct {
	client *Client
}

type Country struct {
	Name                string
	CapitalCity         string
	Iso2Code            string
	Longitude, Latitude string
	Region              Region
	IncomeLevel         IncomeLevel
	LendingType         LendingType
	AdminRegion         struct {
		Id    string
		Value string
	}
}

type Region struct {
	Id    string
	Value string
}

type CountryHeader struct {
	page    int
	pages   int
	perpage string
	total   int
}

type CountryListOptions struct {
	Page    int
	PerPage int
}

// ListCountries list all the countries
func (c *CountryService) ListCountries() ([]Country, error) {
	countryHeader := CountryHeader{}
	country := []Country{}

	req, err := c.client.NewRequest("GET", "countries", nil)
	if err != nil {
		return nil, err
	}

	_, err = c.client.Do(req, &[]interface{}{&countryHeader, &country})

	return country, err
}

// GetCountry Return entity data for particular country
// with ISO 3166-1 alpha-2 code “BR” (ISO 3166-1 alpha-3 code "BRA" can be used in place of alpha-2 code)
func (c *CountryService) GetCountry(countryId string) ([]Country, error) {
	countryHeader := CountryHeader{}
	country := []Country{}

	u := fmt.Sprintf("countries/%v", countryId)
	req, err := c.client.NewRequest("GET", u, nil)

	if err != nil {
		return nil, err
	}
	_, err = c.client.Do(req, &[]interface{}{&countryHeader, &country})

	return country, err

}
