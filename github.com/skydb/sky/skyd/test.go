package skyd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

func assertProperty(t *testing.T, property *Property, id int64, name string, transient bool, dataType string) {
	if property.Id != id {
		t.Fatalf("Unexpected property id. Expected %v, got %v", id, property.Id)
	}
	if property.Name != name {
		t.Fatalf("Unexpected property name. Expected %v, got %v", name, property.Name)
	}
	if property.Transient != transient {
		t.Fatalf("Unexpected property transiency. Expected %v, got %v", transient, property.Transient)
	}
	if property.DataType != dataType {
		t.Fatalf("Unexpected property data type. Expected %v, got %v", dataType, property.DataType)
	}
}

func assertResponse(t *testing.T, resp *http.Response, statusCode int, content string, message string) {
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 || content != string(body) {
		t.Fatalf("%v:\nexp:[%v] %s\ngot:[%v] %s.", message, statusCode, content, resp.StatusCode, string(body))
	}
}

func sendTestHttpRequest(method string, url string, contentType string, body string) (*http.Response, error) {
	client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
	req, _ := http.NewRequest(method, url, strings.NewReader(body))
	req.Header.Add("Content-Type", contentType)
	return client.Do(req)
}

func runTestServer(f func(s *Server)) {
	path, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(path)
	server := NewServer(8586, path)
	server.Silence()
	server.ListenAndServe(nil)
	defer server.Shutdown()
	f(server)
}

func createTempTable(t *testing.T) *Table {
	path, err := ioutil.TempDir("", "")
	os.RemoveAll(path)

	table := NewTable("test", path)
	err = table.Create()
	if err != nil {
		t.Fatalf("Unable to create table: %v", err)
	}

	return table
}

func setupTestTable(name string) {
	resp, _ := sendTestHttpRequest("POST", "http://localhost:8586/tables", "application/json", fmt.Sprintf(`{"name":"%v"}`, name))
	resp.Body.Close()
}

func setupTestProperty(tableName string, name string, transient bool, dataType string) {
	resp, _ := sendTestHttpRequest("POST", fmt.Sprintf("http://localhost:8586/tables/%v/properties", tableName), "application/json", fmt.Sprintf(`{"name":"%v", "transient":%v, "dataType":"%v"}`, name, transient, dataType))
	resp.Body.Close()
}

func setupTestData(t *testing.T, tableName string, items [][]string) {
	for i, item := range items {
		resp, _ := sendTestHttpRequest("PUT", fmt.Sprintf("http://localhost:8586/tables/%s/objects/%s/events/%s", tableName, item[0], item[1]), "application/json", item[2])
		resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Fatalf("setupTestData[%d]: Expected 200, got %v.", i, resp.StatusCode)
		}
	}
}

func _codegen(t *testing.T, tableName string, query string) {
	resp, _ := sendTestHttpRequest("POST", fmt.Sprintf("http://localhost:8586/tables/%s/query/codegen", tableName), "application/json", query)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
