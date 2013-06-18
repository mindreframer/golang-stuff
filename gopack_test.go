package main

import (
	"fmt"
	"testing"
)

func TestUnusedDep(t *testing.T) {
	errors := findErrors(fmt.Sprintf("%s/unused-dep", GopackTestProjects), t)
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, found %d\n", len(errors))
	}
	e := errors[0]
	if e.Kind != UnusedDep {
		t.Errorf("expected unused dependency error\n")
	}
}

func TestUnmanagedImport(t *testing.T) {
	errors := findErrors(fmt.Sprintf("%s/unmanaged-import", GopackTestProjects), t)
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, found %d\n", len(errors))
	}
	e := errors[0]
	if e.Kind != UnmanagedImport {
		t.Errorf("expected unmanaged import error\n")
	}
}

func findErrors(dir string, t *testing.T) []*ProjectError {
	d := LoadDependencyModel(dir)
	p, err := AnalyzeSourceTree(dir)
	if err != nil {
		t.Fatal(err)
	}
	errors := d.Validate(p)
	PrintErrors(errors, t)
	return errors
}

func PrintErrors(errors []*ProjectError, t *testing.T) {
	for _, e := range errors {
		t.Logf("%s\n", e.String())
	}
}
