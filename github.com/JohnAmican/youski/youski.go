package youski

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "net/url"
  
  "appengine"
  "appengine/urlfetch"
)

const (
  BaseURI string = "https://www.googleapis.com/youtube/v3"
  SearchPath string = "/search?key=AIzaSyCYkW3D7NHYpx9TLBzHdh37YXYlUUBt860&part=snippet&type=video"
)

type Results struct {
  Info PageInfo `json:"pageInfo"`
  Items []Item `json:"items"`
}

type PageInfo struct {
  TotalResults int `json:"totalResults"`
  ResultsPerPage int `json:"resultsPerPage"`
}

type Item struct {
  Id ItemId `json:”id”`
  Snippet Snippet `json:”snippet”`
}

type ItemId struct {
  VideoId string `json:”videoId”`
}

type Snippet struct {
  Title string `json:”title”`
}

type ListEntry struct {
  Title string
  VideoId string
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
  body, _ := ioutil.ReadFile("templates/welcome.html")
  fmt.Fprint(w, string(body))
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
  result := search(r, "query")
  respondWith(w, result)
}

func relatedHandler(w http.ResponseWriter, r *http.Request) {
  result := related(r, "videoId")
  respondWith(w, result)
}

func respondWith(w http.ResponseWriter, val interface{}) {
  w.Header().Set("Content-Type", "application/json")
  resp, _ := json.Marshal(val)
  fmt.Fprint(w, string(resp))
}

func init() {
  http.HandleFunc("/", welcomeHandler)
  http.HandleFunc("/search", searchHandler)
  http.HandleFunc("/related", relatedHandler)
  
  http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
  // http.ListenAndServe(":8080", nil)
}

func search(r *http.Request, formValue string) []ListEntry {
  url := BaseURI + SearchPath + "&q=" + url.QueryEscape(r.FormValue(formValue))
  return getListEntries(r, url);
}

func related(r *http.Request, formValue string) []ListEntry {
  url := BaseURI + SearchPath + "&maxResults=50&relatedToVideoId=" + r.FormValue(formValue)
  return getListEntries(r, url);
}

func getListEntries(r *http.Request, url string) []ListEntry {
  c := appengine.NewContext(r)
  client := urlfetch.Client(c)
  resp, err := client.Get(url)
  log.Print("hello!")
  log.Print(resp)
  log.Print(err)
  log.Print("OK!")
  
  defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body)
  
  var results Results
  _ = json.Unmarshal(body, &results)
  
  it := make([]ListEntry, results.length())
  for i, result := range results.Items {
    it[i] = ListEntry{result.Snippet.Title, result.Id.VideoId}
  }
  return it
}

func (result *Results) length() int {
  info := result.Info
  
  if info.TotalResults < info.ResultsPerPage {
    return info.TotalResults
  } else {
    return info.ResultsPerPage
  }
}
