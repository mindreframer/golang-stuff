package wbdata

import (
	"fmt"
)

// TopicsService provides access to the topic sources queries
// in the World Bank Open Data API.
//
// World Bank Open Data API docs: http://data.worldbank.org/node/209
type TopicsService struct {
	client *Client
}

type TopicHeader struct {
	Page    int
	Pages   int
	PerPage int
	Total   int
}

type Topic struct {
	Id, Value, SourceNote string
}

// LisTopics retrieves topics which are high level categories that all indicators are mapped to.
func (t *TopicsService) ListTopics() ([]Topic, error) {

	header := TopicHeader{}
	topic := []Topic{}

	req, err := t.client.NewRequest("GET", "topics", nil)
	if err != nil {
		return nil, err
	}

	_, err = t.client.Do(req, &[]interface{}{&header, &topic})

	return topic, err

}

func (t *TopicsService) GetTopic(id int) ([]Topic, error) {
	header := TopicHeader{}
	topic := []Topic{}

	urlStr := fmt.Sprintf("topic/%v", id)

	req, err := t.client.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	_, err = t.client.Do(req, &[]interface{}{&header, &topic})

	return topic, err

}
