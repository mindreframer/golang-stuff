package dynamodb

import (
	simplejson "github.com/bitly/go-simplejson"
)

func (t *Table) DescribeTable() (*TableDescriptionT, error) {
	q := NewQuery(t)
	jsonResponse, err := t.Server.queryServer(target("DescribeTable"), q)
	if err != nil {
		return nil, err
	}

	json, err := simplejson.NewJson(jsonResponse)
	if err != nil {
		return nil, err
	}

	var tableDescription TableDescriptionT

	// TODO: Populate tableDescription.AttributeDefinitions.

	tableDescription.CreationDateTime, _ = json.Get("CreationDateTime").Float64()
	tableDescription.ItemCount, _ = json.Get("ItemCount").Int64()

	// TODO: Populate tableDescription.KeySchema.
	// TODO: Populate tableDescription.LocalSecondaryIndexes.
	// TODO: Populate tableDescription.ProvisionedThroughPut.

	tableDescription.TableName, _ = json.Get("TableName").String()
	tableDescription.TableSizeBytes, _ = json.Get("TableSizeBytes").Int64()
	tableDescription.TableStatus, _ = json.Get("TableStatus").String()

	return &tableDescription, nil
}
