package dynamodb

import (
	"fmt"
	"github.com/crowdmob/goamz/aws"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	Auth   aws.Auth
	Region aws.Region
}

/*
type Query struct {
	Query string
}
*/

/*
func NewQuery(queryParts []string) *Query {
	return &Query{
		"{" + strings.Join(queryParts, ",") + "}",
	}
}
*/

func (s *Server) queryServer(target string, query *Query) ([]byte, error) {
	data := strings.NewReader(query.String())
	hreq, err := http.NewRequest("POST", s.Region.DynamoDBEndpoint+"/", data)
	if err != nil {
		return nil, err
	}

	hreq.Header.Set("Date", requestDate())
	hreq.Header.Set("Content-Type", "application/x-amz-json-1.0")
	hreq.Header.Set("X-Amz-Target", target)

	service := Service{
		"dynamodb",
		s.Region.Name,
	}

	err = service.Sign(&s.Auth, hreq)

	if err == nil {

		resp, err := http.DefaultClient.Do(hreq)

		if err != nil {
			fmt.Printf("Error calling Amazon")
			return nil, err
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			fmt.Printf("Could not read response body")
			return nil, err
		}

		return body, nil

	}

	return nil, err

}

func requestDate() string {
	now := time.Now().UTC()
	return now.Format(http.TimeFormat)
}

func target(name string) string {
	return "DynamoDB_20111205." + name
}
