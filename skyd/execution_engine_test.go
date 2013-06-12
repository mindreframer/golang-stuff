package skyd

import (
	"testing"
)

// Ensure that the lua script can extract event property references.
func TestExecutionEngineExtractPropertyReferences(t *testing.T) {
	table := createTempTable(t)
	table.Open()
	defer table.Close()

	table.CreateProperty("name", false, "string")
	table.CreateProperty("salary", false, "float")
	table.CreateProperty("purchaseAmount", true, "integer")
	table.CreateProperty("isMember", true, "boolean")

	l, err := NewExecutionEngine(table, "function f(event) x = event:name() if event.salary > 100 then print(event.purchaseAmount, event, event:name()) end end")
	if err != nil {
		t.Fatalf("Unable to create execution engine: %v", err)
	}
	if len(l.propertyRefs) != 3 {
		t.Fatalf("Expected %v properties, got %v", 3, len(l.propertyRefs))
	}
	if p, _ := table.GetPropertyByName("purchaseAmount"); p != l.propertyRefs[0] {
		t.Fatalf("Expected %v, got %v", p, l.propertyRefs[0])
	}
	if p, _ := table.GetPropertyByName("name"); p != l.propertyRefs[1] {
		t.Fatalf("Expected %v, got %v", p, l.propertyRefs[1])
	}
	if p, _ := table.GetPropertyByName("salary"); p != l.propertyRefs[2] {
		t.Fatalf("Expected %v, got %v", p, l.propertyRefs[2])
	}
}
