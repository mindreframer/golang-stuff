package dynamodb_test

import simplejson "github.com/bitly/go-simplejson"
import (
	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/dynamodb"
	"testing"
)

func TestEmptyQuery(t *testing.T) {
	q := dynamodb.NewEmptyQuery()
	queryString := q.String()
	expectedString := "{}"

	if expectedString != queryString {
		t.Fatalf("Unexpected Query String : %s\n", queryString)
	}

}

func TestGetItemQuery(t *testing.T) {
	auth := &aws.Auth{"", "wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY"}
	server := dynamodb.Server{*auth, aws.USEast}
	primary := dynamodb.NewStringAttribute("domain", "")
	key := dynamodb.PrimaryKey{primary, nil}
	table := server.NewTable("sites", key)

	q := dynamodb.NewQuery(table)
	q.AddKey(table, "test", "")

	queryString := []byte(q.String())

	json, err := simplejson.NewJson(queryString)

	if err != nil {
		t.Logf("JSON err : %s\n", err)
		t.Fatalf("Invalid JSON : %s\n", queryString)
	}

	tableName := json.Get("TableName").MustString()

	if tableName != "sites" {
		t.Fatalf("Expected tableName to be sites was : %s", tableName)
	}

	keyMap, err := json.Get("Key").Map()

	if err != nil {
		t.Fatalf("Expected a Key")
	}

	hashRangeKey := keyMap["HashKeyElement"]

	if hashRangeKey == nil {
		t.Fatalf("Expected a HashKeyElement found : %s", keyMap)
	}

	if v, ok := hashRangeKey.(map[string]interface{}); ok {
		if val, ok := v["S"].(string); ok {
			if val != "test" {
				t.Fatalf("Expected HashKeyElement to have the value 'test' found : %s", val)
			}
		}
	} else {
		t.Fatalf("HashRangeKeyt had the wrong type found : %s", hashRangeKey)
	}
}

func TestUpdateQuery(t *testing.T) {
	auth := &aws.Auth{"", "wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY"}
	server := dynamodb.Server{*auth, aws.USEast}
	primary := dynamodb.NewStringAttribute("domain", "")
	rangek := dynamodb.NewNumericAttribute("time", "")
	key := dynamodb.PrimaryKey{primary, rangek}
	table := server.NewTable("sites", key)

	countAttribute := dynamodb.NewNumericAttribute("count", "4")
	attributes := []dynamodb.Attribute{*countAttribute}

	q := dynamodb.NewQuery(table)
	q.AddKey(table, "test", "1234")
	q.AddUpdates(attributes, "ADD")

	queryString := []byte(q.String())

	json, err := simplejson.NewJson(queryString)

	if err != nil {
		t.Logf("JSON err : %s\n", err)
		t.Fatalf("Invalid JSON : %s\n", queryString)
	}

	tableName := json.Get("TableName").MustString()

	if tableName != "sites" {
		t.Fatalf("Expected tableName to be sites was : %s", tableName)
	}

	keyMap, err := json.Get("Key").Map()

	if err != nil {
		t.Fatalf("Expected a Key")
	}

	hashRangeKey := keyMap["HashKeyElement"]

	if hashRangeKey == nil {
		t.Fatalf("Expected a HashKeyElement found : %s", keyMap)
	}

	rangeKey := keyMap["RangeKeyElement"]

	if rangeKey == nil {
		t.Fatalf("Expected a RangeKeyElement found : %s", keyMap)
	}

}
