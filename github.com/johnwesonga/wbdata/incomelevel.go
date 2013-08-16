package wbdata

// IncomeLevelService show the income category of a
// particular country as identified by the World Bank.
//
// World Bank Open Data API docs: http://data.worldbank.org/node/207
type IncomeLevelService struct {
	client *Client
}

type IncomeLevel struct {
	Id    string
	Value string
}

type IncomeLevelHeader struct {
	page    string
	pages   string
	perpage string
	total   string
}

func (i *IncomeLevelService) ListIncomeLevels() ([]IncomeLevel, error) {
	incomeLevelHeader := IncomeLevelHeader{}
	incomeLevel := []IncomeLevel{}

	req, err := i.client.NewRequest("GET", "incomeLevels", nil)
	if err != nil {
		return nil, err
	}

	_, err = i.client.Do(req, &[]interface{}{&incomeLevelHeader, &incomeLevel})

	return incomeLevel, err

}
