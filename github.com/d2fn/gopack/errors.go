package main

import (
	"fmt"
)

const (
	UnusedDep       = "unused-dep"
	UnmanagedImport = "unmanaged-import"
)

type ProjectError struct {
	Kind    string
	Message string
}

func UnusedDependencyError(importPath string) *ProjectError {
	return &ProjectError{
		UnusedDep,
		fmt.Sprintf("%s in gopack.config is unused", importPath),
	}
}

func UnmanagedImportError(s *ImportStats) *ProjectError {
	msg := fmt.Sprintf("%s referenced in the following locations but not managed in gopack.config\n%s", s.Path, s.ReferenceList())
	return &ProjectError{
		UnmanagedImport,
		msg,
	}
}

func (e *ProjectError) String() string {
	return e.Message
}

func (e *ProjectError) Error() string {
	return e.String()
}
