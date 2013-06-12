// +build windows

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
	// Because of our test structure, this method is called several times 
	// during a test run (but it should only be called once during the lifetime
	// of the program).  Windows truly hates this, so we comment it out for it.
	// Other OSes don't seem to care.
	//C.xmlCleanupParser()
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
