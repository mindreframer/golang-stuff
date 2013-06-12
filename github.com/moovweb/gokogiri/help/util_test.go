package help

import "testing"

func CheckXmlMemoryLeaks(t *testing.T) {
	LibxmlCleanUpParser()
	if ! LibxmlCheckMemoryLeak() {
		t.Errorf("Memory leaks: %d!!!", LibxmlGetMemoryAllocation())
		LibxmlReportMemoryLeak()
	}
}
