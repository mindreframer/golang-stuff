package xpath

/*
#cgo CFLAGS: -I../../../clibs/include/libxml2
#cgo LDFLAGS: -lxml2 -L../../../clibs/lib
#include <libxml/xpath.h>
#include <libxml/xpathInternals.h>
#include <libxml/parser.h>

xmlNode* fetchNode(xmlNodeSet *nodeset, int index) {
  	return nodeset->nodeTab[index];
}
*/
import "C"

import "time"
import "unsafe"
import . "gokogiri/util"
import "runtime"
import "errors"

type XPath struct {
	ContextPtr *C.xmlXPathContext
	ResultPtr  *C.xmlXPathObject
}

func NewXPath(docPtr unsafe.Pointer) (xpath *XPath) {
	if docPtr == nil {
		return
	}
	xpath = &XPath{ContextPtr: C.xmlXPathNewContext((*C.xmlDoc)(docPtr)), ResultPtr: nil}
	runtime.SetFinalizer(xpath, (*XPath).Free)
	return
}

func (xpath *XPath) RegisterNamespace(prefix, href string) bool {
	var prefixPtr unsafe.Pointer = nil
	if len(prefix) > 0 {
		prefixBytes := AppendCStringTerminator([]byte(prefix))
		prefixPtr = unsafe.Pointer(&prefixBytes[0])
	}

	var hrefPtr unsafe.Pointer = nil
	if len(href) > 0 {
		hrefBytes := AppendCStringTerminator([]byte(href))
		hrefPtr = unsafe.Pointer(&hrefBytes[0])
	}

	result := C.xmlXPathRegisterNs(xpath.ContextPtr, (*C.xmlChar)(prefixPtr), (*C.xmlChar)(hrefPtr))
	return result == 0
}

//need to add an error as a return value b/c xpath evaluation can return error
func (xpath *XPath) Evaluate(nodePtr unsafe.Pointer, xpathExpr *Expression) (nodes []unsafe.Pointer, err error) {
	if nodePtr == nil {
		//evaluating xpath on a  nil node returns no result.
		return
	}
	xpath.ContextPtr.node = (*C.xmlNode)(nodePtr)
	if xpath.ResultPtr != nil {
		C.xmlXPathFreeObject(xpath.ResultPtr)
	}

	xpath.ResultPtr = C.xmlXPathCompiledEval(xpathExpr.Ptr, xpath.ContextPtr)
	if xpath.ResultPtr == nil {
		err = errors.New("err in evaluating xpath: " + xpathExpr.String())
		return
	}

	if nodesetPtr := xpath.ResultPtr.nodesetval; nodesetPtr != nil {
		if nodesetSize := int(nodesetPtr.nodeNr); nodesetSize > 0 {
			nodes = make([]unsafe.Pointer, nodesetSize)
			for i := 0; i < nodesetSize; i++ {
				nodes[i] = unsafe.Pointer(C.fetchNode(nodesetPtr, C.int(i)))
			}
		}
	}
	return
}

func (xpath *XPath) SetDeadline(deadline *time.Time) {
	if deadline == nil {
		C.xmlXPathContextSetDeadline(xpath.ContextPtr, C.time_t(0))
	} else {
		t := deadline.Unix()
		C.xmlXPathContextSetDeadline(xpath.ContextPtr, C.time_t(t))
	}
}

func (xpath *XPath) Free() {
	if xpath.ContextPtr != nil {
		C.xmlXPathFreeContext(xpath.ContextPtr)
		xpath.ContextPtr = nil
	}
	if xpath.ResultPtr != nil {
		C.xmlXPathFreeObject(xpath.ResultPtr)
		xpath.ResultPtr = nil
	}
}
