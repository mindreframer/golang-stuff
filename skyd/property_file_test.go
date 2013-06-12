package skyd

import (
	"bytes"
	"testing"
)

// Encode a property file.
func TestPropertyFileEncode(t *testing.T) {
	p := NewPropertyFile("")
	p.CreateProperty("name", false, "string")
	p.CreateProperty("salary", false, "float")
	p.CreateProperty("purchaseAmount", true, "integer")
	p.CreateProperty("isMember", true, "boolean")

	// Encode
	buffer := new(bytes.Buffer)
	err := p.Encode(buffer)
	if err != nil {
		t.Fatalf("Unable to encode property file: %v", err)
	}
	expected := `[{"id":-2,"name":"isMember","transient":true,"dataType":"boolean"},{"id":-1,"name":"purchaseAmount","transient":true,"dataType":"integer"},{"id":1,"name":"name","transient":false,"dataType":"string"},{"id":2,"name":"salary","transient":false,"dataType":"float"}]` + "\n"
	if buffer.String() != expected {
		t.Fatalf("Invalid property file encoding:\nexp: %v\ngot: %v", expected, buffer.String())
	}
}

// Decode a property file.
func TestPropertyFileDecode(t *testing.T) {
	p := NewPropertyFile("")

	// Decode
	buffer := bytes.NewBufferString("[{\"id\":-1,\"name\":\"purchaseAmount\",\"transient\":true,\"dataType\":\"integer\"},{\"id\":1,\"name\":\"name\",\"transient\":false,\"dataType\":\"string\"},{\"id\":2,\"name\":\"salary\",\"transient\":false,\"dataType\":\"float\"}, {\"id\":-2,\"name\":\"isMember\",\"transient\":true,\"dataType\":\"boolean\"}]\n")
	err := p.Decode(buffer)
	if err != nil {
		t.Fatalf("Unable to decode property file: %v", err)
	}
	assertProperty(t, p.properties[-2], -2, "isMember", true, "boolean")
	assertProperty(t, p.properties[-1], -1, "purchaseAmount", true, "integer")
	assertProperty(t, p.properties[1], 1, "name", false, "string")
	assertProperty(t, p.properties[2], 2, "salary", false, "float")
}

// Convert a map of string keys into property id keys.
func TestPropertyFileNormalizeMap(t *testing.T) {
	p := NewPropertyFile("")
	p.CreateProperty("name", false, "string")
	p.CreateProperty("salary", false, "float")
	p.CreateProperty("purchaseAmount", true, "integer")

	m := map[string]interface{}{"name": "bob", "salary": 100, "purchaseAmount": 12}
	ret, err := p.NormalizeMap(m)
	if err != nil {
		t.Fatalf("Unable to normalize map: %v", err)
	}
	if ret[1] != "bob" {
		t.Fatalf("ret[1]: Expected %q, got %q", "bob", ret[1])
	}
	if ret[2] != 100 {
		t.Fatalf("ret[2]: Expected %q, got %q", 100, ret[2])
	}
	if ret[-1] != 12 {
		t.Fatalf("ret[-1]: Expected %q, got %q", 12, ret[-1])
	}
}

// Convert a map of string keys into property id keys.
func TestPropertyFileDenormalizeMap(t *testing.T) {
	p := NewPropertyFile("")
	p.CreateProperty("name", false, "string")
	p.CreateProperty("salary", false, "float")
	p.CreateProperty("purchaseAmount", true, "integer")

	m := map[int64]interface{}{1: "bob", 2: 100, -1: 12}
	ret, err := p.DenormalizeMap(m)
	if err != nil {
		t.Fatalf("Unable to denormalize map: %v", err)
	}
	if ret["name"] != "bob" {
		t.Fatalf("ret[\"name\"]: Expected %q, got %q", "bob", ret["name"])
	}
	if ret["salary"] != 100 {
		t.Fatalf("ret[\"salary\"]: Expected %q, got %q", 100, ret["salary"])
	}
	if ret["purchaseAmount"] != 12 {
		t.Fatalf("ret[\"purchaseAmount\"]: Expected %q, got %q", 12, ret["purchaseAmount"])
	}
}
