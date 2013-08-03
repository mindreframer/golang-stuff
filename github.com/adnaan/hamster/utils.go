/*Common utility functions and methods*/
package hamster

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	//"os"
	"errors"
	"strings"
	//"testing"
)

var (
	port        = 8686
	host        = "http://localhost:8686"
	mongoHost   = "mongodb://adnaan:pass@localhost:27017/hamster"
	contentType = "application/json"
)

func testHttpRequest(verb string, resource string, body string) (*http.Response, error) {
	client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
	r, _ := http.NewRequest(verb, fmt.Sprintf("%s%s", host, resource), strings.NewReader(body))
	r.Header.Add("Content-Type", contentType)
	return client.Do(r)
}
func testHttpRequestWithHeaders(verb string, resource string, body string, header map[string]string) (*http.Response, error) {

	client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
	r, _ := http.NewRequest(verb, fmt.Sprintf("%s%s", host, resource), strings.NewReader(body))
	r.Header.Add("Content-Type", contentType)
	for key, value := range header {
		r.Header.Add(key, value)
	}
	return client.Do(r)

}

func testPostPng(resource string, fileReader io.Reader, header map[string]string) (*http.Response, error) {
	client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
	r, _ := http.NewRequest("POST", fmt.Sprintf("%s%s", host, resource), fileReader)
	r.Header.Add("Content-Type", "image/png")
	for key, value := range header {
		r.Header.Add(key, value)
	}
	return client.Do(r)
}

func testHttp(verb string, resource string, header map[string]string) (*http.Response, error) {
	client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
	r, _ := http.NewRequest(verb, fmt.Sprintf("%s%s", host, resource), nil)

	for key, value := range header {
		r.Header.Add(key, value)
	}
	return client.Do(r)

}

func testServer(f func(s *Server)) {

	server := NewServer()
	//server.Quiet()
	server.ListenAndServe()
	defer server.Shutdown()
	f(server)
}

func (s *Server) readJson(d interface{}, r *http.Request, w http.ResponseWriter) error {

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {

		s.logger.Printf("error in reading body for: %v, err: %v\n ", r.Body, err)
		http.Error(w, "Bad Data!", http.StatusBadRequest)
		return err
	}

	return json.Unmarshal(body, &d)

}

func (s *Server) serveJson(w http.ResponseWriter, v interface{}) {
	content, err := json.MarshalIndent(v, "", "  ")
	if err != nil {

		s.logger.Printf("error in serving json err: %v  \n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.Header().Set("Content-Type", "application/json")
	w.Write(content)

}

func (s *Server) getObjectId(w http.ResponseWriter, r *http.Request) string {
	oid := r.URL.Query().Get(":objectId")
	if oid == "" {
		s.notFound(r, w, errors.New("objectId is empty"), "val: "+oid)
	}

	object_id := decodeToken(oid)
	if object_id == "" {
		s.notFound(r, w, errors.New("objectId cannot be decoded"), "val: "+oid)
	}

	return object_id

}

func (s *Server) getObjectName(w http.ResponseWriter, r *http.Request) string {
	oname := r.URL.Query().Get(":objectName")
	if oname == "" {
		s.notFound(r, w, errors.New("objectName is empty"), "val: "+oname)
	}

	return oname

}

func (s *Server) getFileName(w http.ResponseWriter, r *http.Request) string {
	fname := r.URL.Query().Get(":fileName")
	if fname == "" {
		s.notFound(r, w, errors.New("fileName is empty"), "val: "+fname)
	}

	return fname

}

func (s *Server) getFileParams(w http.ResponseWriter, r *http.Request) (string, string) {
	fname := r.URL.Query().Get(":fileName")
	if fname == "" {
		s.notFound(r, w, errors.New("fileName is empty"), "val: "+fname)
	}

	fid := r.URL.Query().Get(":fileId")
	if fid == "" {
		s.notFound(r, w, errors.New("fileId is empty"), "val: "+fid)
	}

	file_id := decodeToken(fid)
	if file_id == "" {
		s.notFound(r, w, errors.New("object params cannot be decoded"), "val: "+fname+" , "+file_id)
	}

	return fname, file_id

}

func (s *Server) getObjectParams(w http.ResponseWriter, r *http.Request) (string, string) {
	object_name := r.URL.Query().Get(":objectName")
	oid := r.URL.Query().Get(":objectId")
	if oid == "" || object_name == "" {
		s.notFound(r, w, errors.New("object params are invalid"), "val: "+object_name+" , "+oid)
	}

	object_id := decodeToken(oid)
	if object_id == "" {
		s.notFound(r, w, errors.New("object params cannot be decoded"), "val: "+object_name+" , "+oid)
	}

	return object_name, object_id

}

func (s *Server) getAppObjectId(w http.ResponseWriter, r *http.Request) string {
	atok := r.Header.Get("X-Api-Token")

	if atok == "" {
		s.unauthorized(r, w, errors.New("token is empty"), "api token invalid")
	}

	object_id := decodeToken(atok)
	if object_id == "" {
		s.notFound(r, w, errors.New("app objectid cannot be decoded"), "val: "+atok)
	}

	return object_id

}

func (s *Server) getAppParams(w http.ResponseWriter, r *http.Request) (string, string) {
	did := r.URL.Query().Get(":developerId")
	oid := r.URL.Query().Get(":objectId")
	if oid == "" || did == "" {
		s.notFound(r, w, errors.New("app params are invalid"), "val: "+did+" , "+oid)
	}

	developer_id := decodeToken(did)
	object_id := decodeToken(oid)
	if object_id == "" || developer_id == "" {
		s.notFound(r, w, errors.New("app params cannot be decoded"), "val: "+did+" , "+oid)
	}

	return developer_id, object_id

}
