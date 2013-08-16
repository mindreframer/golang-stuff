package wbdata

// LendingTypeService provides access to the countries related queries
// in the World Bank Open Data API.
//
// World Bank Open Data API docs: http://data.worldbank.org/node/208
type LendingTypeService struct {
	client *Client
}

type LendingType struct {
	Id    string
	Value string
}

type LendingTypeHeader struct {
	page    int
	pages   int
	perpage string
	total   int
}

func (l *LendingTypeService) ListLendingTypes() ([]LendingType, error) {
	lendingTypeHeader := LendingTypeHeader{}
	lendingType := []LendingType{}

	req, err := l.client.NewRequest("GET", "lendingTypes", nil)
	if err != nil {
		return nil, err
	}

	_, err = l.client.Do(req, &[]interface{}{&lendingTypeHeader, &lendingType})

	return lendingType, err

}
