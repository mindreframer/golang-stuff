// +build !windows

package help

/*
#cgo CFLAGS: -I../../../clibs/include/libxml2
#cgo LDFLAGS: -lxml2 -L../../../clibs/lib

#include <libxml/tree.h>
#include <libxml/parser.h>
#include <libxml/HTMLtree.h>
#include <libxml/HTMLparser.h>
#include <libxml/xmlsave.h>

void printMemoryLeak() { xmlMemDisplay(stdout); }
*/
import "C"

func LibxmlInitParser() {
	C.xmlInitParser()
}

func LibxmlCleanUpParser() {
	C.xmlCleanupParser()
}

func LibxmlGetMemoryAllocation() int {
	return (int)(C.xmlMemBlocks())
}

func LibxmlCheckMemoryLeak() bool {
	return (C.xmlMemBlocks() == 0)
}

func LibxmlReportMemoryLeak() {
	C.printMemoryLeak()
}
