package wbdata

// SourcesService provides access to the catalog sources queries
// in the World Bank Open Data API.
//
// World Bank Open Data API docs: http://data.worldbank.org/node/210
type SourcesService struct {
	client *Client
}

type SourceHeader struct {
	Page    string
	Pages   string
	PerPage string
	Total   string
}

type Source struct {
	Id          string
	Name        string
	Description string
	Url         string
}

func (s *SourcesService) ListSources() ([]Source, error) {
	sourceHeader := SourceHeader{}
	source := []Source{}

	req, err := s.client.NewRequest("GET", "sources", nil)
	if err != nil {
		return nil, err
	}

	_, err = s.client.Do(req, &[]interface{}{&sourceHeader, &source})

	return source, err

}
