package skyd

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

// Ensure that we can create a new table.
func TestFactorization(t *testing.T) {
	path, err := ioutil.TempDir("", "")
	defer os.RemoveAll(path)
	path = fmt.Sprintf("%v/factors", path)

	factors := NewFactors(path)
	defer factors.Close()
	err = factors.Open()
	if err != nil {
		t.Fatalf("Unable to create factors: %v", err)
	}

	num, err := factors.Factorize("foo", "bar", "/index.html", true)
	if err != nil || num != 1 {
		t.Fatalf("Wrong factorization: exp: %v, got: %v (%v)", 1, num, err)
	}
	num, err = factors.Factorize("foo", "bar", "/about.html", true)
	if err != nil || num != 2 {
		t.Fatalf("Wrong factorization: exp: %v, got: %v (%v)", 2, num, err)
	}

	str, err := factors.Defactorize("foo", "bar", 1)
	if err != nil || str != "/index.html" {
		t.Fatalf("Wrong defactorization: exp: %v, got: %v (%v)", "/index.html", str, err)
	}
	str, err = factors.Defactorize("foo", "bar", 2)
	if err != nil || str != "/about.html" {
		t.Fatalf("Wrong defactorization: exp: %v, got: %v (%v)", "/about.html", str, err)
	}
}
