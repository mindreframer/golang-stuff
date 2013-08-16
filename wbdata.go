/*
Package wbdata provides a client for using the World Bank Open Data API.

Access different parts of the World Bank Open Data API using the various
services:
         client := wbdata.NewClient(nil)

         // list all countries
         countries, err := client.Countries.GetCountries()


The full World Bank Open Data API is documented at http://data.worldbank.org/developers/api-overview.
*/
package wbdata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	libraryVersion = "0.1"
	defaultBaseURL = "http://api.worldbank.org"
)

// A Client manages communication with the World Bank Open Data API.
type Client struct {
	client *http.Client

	BaseURL *url.URL

	//Services to talk to different APIs
	Countries    *CountryService
	Sources      *SourcesService
	Topics       *TopicsService
	Indicators   *IndicatorService
	IncomeLevels *IncomeLevelService
	LendingTypes *LendingTypeService
}

func NewClient() *Client {
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{client: http.DefaultClient, BaseURL: baseURL}
	c.Countries = &CountryService{client: c}
	c.Sources = &SourcesService{client: c}
	c.Topics = &TopicsService{client: c}
	c.Indicators = &IndicatorService{client: c}
	c.IncomeLevels = &IncomeLevelService{client: c}
	c.LendingTypes = &LendingTypeService{client: c}
	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.  If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	v := url.Values{}
	v.Set("format", "json")

	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	// the API expects URL+resource?format=json
	url := fmt.Sprintf("%s?%s", u, v.Encode())
	log.Println(url)

	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// Do sends an API request and returns the API response.  The API response is
// decoded and stored in the value pointed to by v, or returned as an error if
// an API error has occurred.
func (c *Client) Do(req *http.Request, v *[]interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}

	return resp, err
}

// ErrorResponse from the API.
// {"message":
//  [{"id":"120","key":"Parameter 'country' has an invalid value","value":"The provided parameter value is not valid"}]]
type ErrorResponse struct {
	Message []struct {
		Id    string
		Key   string
		Value string
	}
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%+v", r.Message)
}

// CheckResponse checks the API response for errors, and returns them if
// present.  A response is considered an error if it has a status code outside
// the 200 range.  API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse.  Any other
// response body will be silently ignored.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse

}
