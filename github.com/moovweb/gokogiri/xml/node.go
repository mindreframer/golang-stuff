package xml

//#include "helper.h"
//#include <string.h>
import "C"

import "time"

import (
	"errors"
	. "gokogiri/util"
	"gokogiri/xpath"
	"strconv"
	"unsafe"
)

var (
	ERR_UNDEFINED_COERCE_PARAM               = errors.New("unexpected parameter type in coerce")
	ERR_UNDEFINED_SET_CONTENT_PARAM          = errors.New("unexpected parameter type in SetContent")
	ERR_UNDEFINED_SEARCH_PARAM               = errors.New("unexpected parameter type in Search")
	ERR_CANNOT_MAKE_DUCMENT_AS_CHILD         = errors.New("cannot add a document node as a child")
	ERR_CANNOT_COPY_TEXT_NODE_WHEN_ADD_CHILD = errors.New("cannot copy a text node when adding it")
)

//xmlNode types
const (
	XML_ELEMENT_NODE       = 1
	XML_ATTRIBUTE_NODE     = 2
	XML_TEXT_NODE          = 3
	XML_CDATA_SECTION_NODE = 4
	XML_ENTITY_REF_NODE    = 5
	XML_ENTITY_NODE        = 6
	XML_PI_NODE            = 7
	XML_COMMENT_NODE       = 8
	XML_DOCUMENT_NODE      = 9
	XML_DOCUMENT_TYPE_NODE = 10
	XML_DOCUMENT_FRAG_NODE = 11
	XML_NOTATION_NODE      = 12
	XML_HTML_DOCUMENT_NODE = 13
	XML_DTD_NODE           = 14
	XML_ELEMENT_DECL       = 15
	XML_ATTRIBUTE_DECL     = 16
	XML_ENTITY_DECL        = 17
	XML_NAMESPACE_DECL     = 18
	XML_XINCLUDE_START     = 19
	XML_XINCLUDE_END       = 20
	XML_DOCB_DOCUMENT_NODE = 21
)

const (
	XML_SAVE_FORMAT   = 1   // format save output
	XML_SAVE_NO_DECL  = 2   //drop the xml declaration
	XML_SAVE_NO_EMPTY = 4   //no empty tags
	XML_SAVE_NO_XHTML = 8   //disable XHTML1 specific rules
	XML_SAVE_XHTML    = 16  //force XHTML1 specific rules
	XML_SAVE_AS_XML   = 32  //force XML serialization on HTML doc
	XML_SAVE_AS_HTML  = 64  //force HTML serialization on XML doc
	XML_SAVE_WSNONSIG = 128 //format with non-significant whitespace
)

type Node interface {
	NodePtr() unsafe.Pointer
	ResetNodePtr()
	MyDocument() Document

	IsValid() bool

	ParseFragment([]byte, []byte, int) (*DocumentFragment, error)

	//
	NodeType() int
	NextSibling() Node
	PreviousSibling() Node

	Parent() Node
	FirstChild() Node
	LastChild() Node
	CountChildren() int
	Attributes() map[string]*AttributeNode

	//
	Coerce(interface{}) ([]Node, error)

	//
	AddChild(interface{}) error
	AddPreviousSibling(interface{}) error
	AddNextSibling(interface{}) error
	InsertBefore(interface{}) error
	InsertAfter(interface{}) error
	InsertBegin(interface{}) error
	InsertEnd(interface{}) error
	SetInnerHtml(interface{}) error
	SetChildren(interface{}) error
	Replace(interface{}) error
	Wrap(string) error
	//Swap(interface{}) os.Error
	//
	////
	SetContent(interface{}) error

	//
	Name() string
	SetName(string)

	//
	Attr(string) string
	SetAttr(string, string) string
	Attribute(string) *AttributeNode

	//
	Path() string

	//
	Duplicate(int) Node
	DuplicateTo(Document, int) Node

	Search(interface{}) ([]Node, error)
	SearchByDeadline(interface{}, *time.Time) ([]Node, error)

	//SetParent(Node)
	//IsComment() bool
	//IsCData() bool
	//IsXml() bool
	//IsHtml() bool
	//IsText() bool
	//IsElement() bool
	//IsFragment() bool
	//

	//
	Unlink()
	Remove()
	ResetChildren()
	//Free()
	////
	ToXml([]byte, []byte) ([]byte, int)
	ToHtml([]byte, []byte) ([]byte, int)
	ToBuffer([]byte) []byte
	String() string
	Content() string
	InnerHtml() string

	RecursivelyRemoveNamespaces() error
	SetNamespace(string, string)
	RemoveDefaultNamespace()
}

//run out of memory
var ErrTooLarge = errors.New("Output buffer too large")

//pre-allocate a buffer for serializing the document
const initialOutputBufferSize = 10 //100K

type XmlNode struct {
	Ptr *C.xmlNode
	Document
	valid bool
}

type WriteBuffer struct {
	Node   *XmlNode
	Buffer []byte
	Offset int
}

func NewNode(nodePtr unsafe.Pointer, document Document) (node Node) {
	if nodePtr == nil {
		return nil
	}
	xmlNode := &XmlNode{
		Ptr:      (*C.xmlNode)(nodePtr),
		Document: document,
		valid:    true,
	}
	nodeType := C.getNodeType((*C.xmlNode)(nodePtr))

	switch nodeType {
	default:
		node = xmlNode
	case XML_ATTRIBUTE_NODE:
		node = &AttributeNode{XmlNode: xmlNode}
	case XML_ELEMENT_NODE:
		node = &ElementNode{XmlNode: xmlNode}
	case XML_CDATA_SECTION_NODE:
		node = &CDataNode{XmlNode: xmlNode}
	case XML_TEXT_NODE:
		node = &TextNode{XmlNode: xmlNode}
	}
	return
}

func (xmlNode *XmlNode) coerce(data interface{}) (nodes []Node, err error) {
	switch t := data.(type) {
	default:
		err = ERR_UNDEFINED_COERCE_PARAM
	case []Node:
		nodes = t
	case *DocumentFragment:
		nodes = t.Children()
	case string:
		f, err := xmlNode.MyDocument().ParseFragment([]byte(t), nil, DefaultParseOption)
		if err == nil {
			nodes = f.Children()
		}
	case []byte:
		f, err := xmlNode.MyDocument().ParseFragment(t, nil, DefaultParseOption)
		if err == nil {
			nodes = f.Children()
		}
	}
	return
}

func (xmlNode *XmlNode) Coerce(data interface{}) (nodes []Node, err error) {
	return xmlNode.coerce(data)
}

//
func (xmlNode *XmlNode) AddChild(data interface{}) (err error) {
	switch t := data.(type) {
	default:
		if nodes, err := xmlNode.coerce(data); err == nil {
			for _, node := range nodes {
				if err = xmlNode.addChild(node); err != nil {
					break
				}
			}
		}
	case *DocumentFragment:
		if nodes, err := xmlNode.coerce(data); err == nil {
			for _, node := range nodes {
				if err = xmlNode.addChild(node); err != nil {
					break
				}
			}
		}
	case Node:
		err = xmlNode.addChild(t)
	}
	return
}

func (xmlNode *XmlNode) AddPreviousSibling(data interface{}) (err error) {
	switch t := data.(type) {
	default:
		if nodes, err := xmlNode.coerce(data); err == nil {
			for _, node := range nodes {
				if err = xmlNode.addPreviousSibling(node); err != nil {
					break
				}
			}
		}
	case *DocumentFragment:
		if nodes, err := xmlNode.coerce(data); err == nil {
			for _, node := range nodes {
				if err = xmlNode.addPreviousSibling(node); err != nil {
					break
				}
			}
		}
	case Node:
		err = xmlNode.addPreviousSibling(t)
	}
	return
}

func (xmlNode *XmlNode) AddNextSibling(data interface{}) (err error) {
	switch t := data.(type) {
	default:
		if nodes, err := xmlNode.coerce(data); err == nil {
			for i := len(nodes) - 1; i >= 0; i-- {
				node := nodes[i]
				if err = xmlNode.addNextSibling(node); err != nil {
					break
				}
			}
		}
	case *DocumentFragment:
		if nodes, err := xmlNode.coerce(data); err == nil {
			for i := len(nodes) - 1; i >= 0; i-- {
				node := nodes[i]
				if err = xmlNode.addNextSibling(node); err != nil {
					break
				}
			}
		}
	case Node:
		err = xmlNode.addNextSibling(t)
	}
	return
}

func (xmlNode *XmlNode) ResetNodePtr() {
	xmlNode.Ptr = nil
	return
}

func (xmlNode *XmlNode) IsValid() bool {
	return xmlNode.valid
}

func (xmlNode *XmlNode) MyDocument() (document Document) {
	document = xmlNode.Document.DocRef()
	return
}

func (xmlNode *XmlNode) NodePtr() (p unsafe.Pointer) {
	p = unsafe.Pointer(xmlNode.Ptr)
	return
}

func (xmlNode *XmlNode) NodeType() (nodeType int) {
	nodeType = int(C.getNodeType(xmlNode.Ptr))
	return
}

func (xmlNode *XmlNode) Path() (path string) {
	pathPtr := C.xmlGetNodePath(xmlNode.Ptr)
	if pathPtr != nil {
		p := (*C.char)(unsafe.Pointer(pathPtr))
		defer C.xmlFreeChars(p)
		path = C.GoString(p)
	}
	return
}

func (xmlNode *XmlNode) NextSibling() Node {
	siblingPtr := (*C.xmlNode)(xmlNode.Ptr.next)
	return NewNode(unsafe.Pointer(siblingPtr), xmlNode.Document)
}

func (xmlNode *XmlNode) PreviousSibling() Node {
	siblingPtr := (*C.xmlNode)(xmlNode.Ptr.prev)
	return NewNode(unsafe.Pointer(siblingPtr), xmlNode.Document)
}

func (xmlNode *XmlNode) CountChildren() int {
	return int(C.xmlLsCountNode(xmlNode.Ptr))
}

func (xmlNode *XmlNode) FirstChild() Node {
	return NewNode(unsafe.Pointer(xmlNode.Ptr.children), xmlNode.Document)
}

func (xmlNode *XmlNode) LastChild() Node {
	return NewNode(unsafe.Pointer(xmlNode.Ptr.last), xmlNode.Document)
}

func (xmlNode *XmlNode) Parent() Node {
	if C.xmlNodePtrCheck(unsafe.Pointer(xmlNode.Ptr.parent)) == C.int(0) {
		return nil
	}
	return NewNode(unsafe.Pointer(xmlNode.Ptr.parent), xmlNode.Document)
}

func (xmlNode *XmlNode) ResetChildren() {
	var p unsafe.Pointer
	for childPtr := xmlNode.Ptr.children; childPtr != nil; {
		nextPtr := childPtr.next
		p = unsafe.Pointer(childPtr)
		C.xmlUnlinkNodeWithCheck((*C.xmlNode)(p))
		xmlNode.Document.AddUnlinkedNode(p)
		childPtr = nextPtr
	}
}

func (xmlNode *XmlNode) SetContent(content interface{}) (err error) {
	switch data := content.(type) {
	default:
		err = ERR_UNDEFINED_SET_CONTENT_PARAM
	case string:
		err = xmlNode.SetContent([]byte(data))
	case []byte:
		contentBytes := GetCString(data)
		contentPtr := unsafe.Pointer(&contentBytes[0])
		C.xmlSetContent(unsafe.Pointer(xmlNode), unsafe.Pointer(xmlNode.Ptr), contentPtr)
	}
	return
}

func (xmlNode *XmlNode) InsertBefore(data interface{}) (err error) {
	err = xmlNode.AddPreviousSibling(data)
	return
}

func (xmlNode *XmlNode) InsertAfter(data interface{}) (err error) {
	err = xmlNode.AddNextSibling(data)
	return
}

func (xmlNode *XmlNode) InsertBegin(data interface{}) (err error) {
	if parent := xmlNode.Parent(); parent != nil {
		if last := parent.LastChild(); last != nil {
			err = last.AddPreviousSibling(data)
		}
	}
	return
}

func (xmlNode *XmlNode) InsertEnd(data interface{}) (err error) {
	if parent := xmlNode.Parent(); parent != nil {
		if first := parent.FirstChild(); first != nil {
			err = first.AddPreviousSibling(data)
		}
	}
	return
}

func (xmlNode *XmlNode) SetChildren(data interface{}) (err error) {
	nodes, err := xmlNode.coerce(data)
	if err != nil {
		return
	}
	xmlNode.ResetChildren()
	err = xmlNode.AddChild(nodes)
	return nil
}

func (xmlNode *XmlNode) SetInnerHtml(data interface{}) (err error) {
	err = xmlNode.SetChildren(data)
	return
}

func (xmlNode *XmlNode) Replace(data interface{}) (err error) {
	err = xmlNode.AddPreviousSibling(data)
	if err != nil {
		return
	}
	xmlNode.Remove()
	return
}

func (xmlNode *XmlNode) Attributes() (attributes map[string]*AttributeNode) {
	attributes = make(map[string]*AttributeNode)
	for prop := xmlNode.Ptr.properties; prop != nil; prop = prop.next {
		if prop.name != nil {
			namePtr := unsafe.Pointer(prop.name)
			name := C.GoString((*C.char)(namePtr))
			attrPtr := unsafe.Pointer(prop)
			attributeNode := NewNode(attrPtr, xmlNode.Document)
			if attr, ok := attributeNode.(*AttributeNode); ok {
				attributes[name] = attr
			}
		}
	}
	return
}

func (xmlNode *XmlNode) Attribute(name string) (attribute *AttributeNode) {
	if xmlNode.NodeType() != XML_ELEMENT_NODE {
		return
	}
	nameBytes := GetCString([]byte(name))
	namePtr := unsafe.Pointer(&nameBytes[0])
	attrPtr := C.xmlHasNsProp(xmlNode.Ptr, (*C.xmlChar)(namePtr), nil)
	if attrPtr == nil {
		return
	} else {
		node := NewNode(unsafe.Pointer(attrPtr), xmlNode.Document)
		if node, ok := node.(*AttributeNode); ok {
			attribute = node
		}
	}
	return
}

func (xmlNode *XmlNode) Attr(name string) (val string) {
	if xmlNode.NodeType() != XML_ELEMENT_NODE {
		return
	}
	nameBytes := GetCString([]byte(name))
	namePtr := unsafe.Pointer(&nameBytes[0])
	valPtr := C.xmlGetProp(xmlNode.Ptr, (*C.xmlChar)(namePtr))
	if valPtr == nil {
		return
	}
	p := unsafe.Pointer(valPtr)
	defer C.xmlFreeChars((*C.char)(p))
	val = C.GoString((*C.char)(p))
	return
}

func (xmlNode *XmlNode) SetAttr(name, value string) (val string) {
	val = value
	if xmlNode.NodeType() != XML_ELEMENT_NODE {
		return
	}
	nameBytes := GetCString([]byte(name))
	namePtr := unsafe.Pointer(&nameBytes[0])

	valueBytes := GetCString([]byte(value))
	valuePtr := unsafe.Pointer(&valueBytes[0])

	C.xmlSetProp(xmlNode.Ptr, (*C.xmlChar)(namePtr), (*C.xmlChar)(valuePtr))
	return
}

func (xmlNode *XmlNode) Search(data interface{}) (result []Node, err error) {
	switch data := data.(type) {
	default:
		err = ERR_UNDEFINED_SEARCH_PARAM
	case string:
		if xpathExpr := xpath.Compile(data); xpathExpr != nil {
			defer xpathExpr.Free()
			result, err = xmlNode.Search(xpathExpr)
		} else {
			err = errors.New("cannot compile xpath: " + data)
		}
	case []byte:
		result, err = xmlNode.Search(string(data))
	case *xpath.Expression:
		xpathCtx := xmlNode.Document.DocXPathCtx()
		nodePtrs, err := xpathCtx.Evaluate(unsafe.Pointer(xmlNode.Ptr), data)
		if nodePtrs == nil || err != nil {
			return nil, err
		}
		for _, nodePtr := range nodePtrs {
			result = append(result, NewNode(nodePtr, xmlNode.Document))
		}
	}
	return
}

func (xmlNode *XmlNode) SearchByDeadline(data interface{}, deadline *time.Time) (result []Node, err error) {
	xpathCtx := xmlNode.Document.DocXPathCtx()
	xpathCtx.SetDeadline(deadline)
	result, err = xmlNode.Search(data)
	xpathCtx.SetDeadline(nil)
	return
}

/*
func (xmlNode *XmlNode) Replace(interface{}) error {

}
func (xmlNode *XmlNode) Swap(interface{}) error {

}
func (xmlNode *XmlNode) SetParent(Node) {

}
func (xmlNode *XmlNode) IsComment() bool {

}
func (xmlNode *XmlNode) IsCData() bool {

}
func (xmlNode *XmlNode) IsXml() bool {

}
func (xmlNode *XmlNode) IsHtml() bool {

}
func (xmlNode *XmlNode) IsText() bool {

}
func (xmlNode *XmlNode) IsElement() bool {

}
func (xmlNode *XmlNode) IsFragment() bool {

}
*/

func (xmlNode *XmlNode) Name() (name string) {
	if xmlNode.Ptr.name != nil {
		p := unsafe.Pointer(xmlNode.Ptr.name)
		name = C.GoString((*C.char)(p))
	}
	return
}

func (xmlNode *XmlNode) SetName(name string) {
	if len(name) > 0 {
		nameBytes := GetCString([]byte(name))
		namePtr := unsafe.Pointer(&nameBytes[0])
		C.xmlNodeSetName(xmlNode.Ptr, (*C.xmlChar)(namePtr))
	}
}

func (xmlNode *XmlNode) Duplicate(level int) Node {
	return xmlNode.DuplicateTo(xmlNode.Document, level)
}

func (xmlNode *XmlNode) DuplicateTo(doc Document, level int) (dup Node) {
	if xmlNode.valid {
		dupPtr := C.xmlDocCopyNode(xmlNode.Ptr, (*C.xmlDoc)(doc.DocPtr()), C.int(level))
		if dupPtr != nil {
			dup = NewNode(unsafe.Pointer(dupPtr), xmlNode.Document)
		}
	}
	return
}

func (xmlNode *XmlNode) serialize(format int, encoding, outputBuffer []byte) ([]byte, int) {
	nodePtr := unsafe.Pointer(xmlNode.Ptr)
	var encodingPtr unsafe.Pointer
	if len(encoding) == 0 {
		encoding = xmlNode.Document.OutputEncoding()
	}
	if len(encoding) > 0 {
		encodingPtr = unsafe.Pointer(&(encoding[0]))
	} else {
		encodingPtr = nil
	}

	wbuffer := &WriteBuffer{Node: xmlNode, Buffer: outputBuffer}
	wbufferPtr := unsafe.Pointer(wbuffer)

	format |= XML_SAVE_FORMAT
	ret := int(C.xmlSaveNode(wbufferPtr, nodePtr, encodingPtr, C.int(format)))
	if ret < 0 {
		panic("output error in xml node serialization: " + strconv.Itoa(ret))
		return nil, 0
	}

	return wbuffer.Buffer, wbuffer.Offset
}

func (xmlNode *XmlNode) ToXml(encoding, outputBuffer []byte) ([]byte, int) {
	return xmlNode.serialize(XML_SAVE_AS_XML, encoding, outputBuffer)
}

func (xmlNode *XmlNode) ToHtml(encoding, outputBuffer []byte) ([]byte, int) {
	return xmlNode.serialize(XML_SAVE_AS_HTML, encoding, outputBuffer)
}

func (xmlNode *XmlNode) ToBuffer(outputBuffer []byte) []byte {
	var b []byte
	var size int
	if docType := xmlNode.Document.DocType(); docType == XML_HTML_DOCUMENT_NODE {
		b, size = xmlNode.ToHtml(nil, outputBuffer)
	} else {
		b, size = xmlNode.ToXml(nil, outputBuffer)
	}
	return b[:size]
}

func (xmlNode *XmlNode) String() string {
	b := xmlNode.ToBuffer(nil)
	if b == nil {
		return ""
	}
	return string(b)
}

func (xmlNode *XmlNode) Content() string {
	contentPtr := C.xmlNodeGetContent(xmlNode.Ptr)
	charPtr := (*C.char)(unsafe.Pointer(contentPtr))
	defer C.xmlFreeChars(charPtr)
	return C.GoString(charPtr)
}

func (xmlNode *XmlNode) InnerHtml() string {
	out := ""

	for child := xmlNode.FirstChild(); child != nil; child = child.NextSibling() {
		out += child.String()
	}
	return out
}

func (xmlNode *XmlNode) Unlink() {
	if int(C.xmlUnlinkNodeWithCheck(xmlNode.Ptr)) != 0 {
		xmlNode.Document.AddUnlinkedNode(unsafe.Pointer(xmlNode.Ptr))
	}
}

func (xmlNode *XmlNode) Remove() {
	if xmlNode.valid && unsafe.Pointer(xmlNode.Ptr) != xmlNode.Document.DocPtr() {
		xmlNode.Unlink()
		xmlNode.valid = false
	}
}

func (xmlNode *XmlNode) addChild(node Node) (err error) {
	nodeType := node.NodeType()
	if nodeType == XML_DOCUMENT_NODE || nodeType == XML_HTML_DOCUMENT_NODE {
		err = ERR_CANNOT_MAKE_DUCMENT_AS_CHILD
		return
	}
	nodePtr := node.NodePtr()
	if xmlNode.NodePtr() == nodePtr {
		return
	}
	ret := xmlNode.isAccestor(nodePtr)
	if ret < 0 {
		return
	} else if ret == 0 {
		if !xmlNode.Document.RemoveUnlinkedNode(nodePtr) {
			C.xmlUnlinkNodeWithCheck((*C.xmlNode)(nodePtr))
		}
		C.xmlAddChild(xmlNode.Ptr, (*C.xmlNode)(nodePtr))
	} else if ret > 0 {
		node.Remove()
	}

	/*
		childPtr := C.xmlAddChild(xmlNode.Ptr, (*C.xmlNode)(nodePtr))
		if nodeType == XML_TEXT_NODE && childPtr != (*C.xmlNode)(nodePtr) {
			//check the retured pointer
			//if it is not the text node just added, it means that the text node is freed because it has merged into other nodes
			//then we should invalid this node, because we do not want to have a dangling pointer
			node.Remove()
		}
	*/
	return
}

func (xmlNode *XmlNode) addPreviousSibling(node Node) (err error) {
	nodeType := node.NodeType()
	if nodeType == XML_DOCUMENT_NODE || nodeType == XML_HTML_DOCUMENT_NODE {
		err = ERR_CANNOT_MAKE_DUCMENT_AS_CHILD
		return
	}
	nodePtr := node.NodePtr()
	if xmlNode.NodePtr() == nodePtr {
		return
	}
	ret := xmlNode.isAccestor(nodePtr)
	if ret < 0 {
		return
	} else if ret == 0 {
		if !xmlNode.Document.RemoveUnlinkedNode(nodePtr) {
			C.xmlUnlinkNodeWithCheck((*C.xmlNode)(nodePtr))
		}
		C.xmlAddPrevSibling(xmlNode.Ptr, (*C.xmlNode)(nodePtr))
	} else if ret > 0 {
		node.Remove()
	}
	/*
		childPtr := C.xmlAddPrevSibling(xmlNode.Ptr, (*C.xmlNode)(nodePtr))
		if nodeType == XML_TEXT_NODE && childPtr != (*C.xmlNode)(nodePtr) {
			//check the retured pointer
			//if it is not the text node just added, it means that the text node is freed because it has merged into other nodes
			//then we should invalid this node, because we do not want to have a dangling pointer
			//xmlNode.Document.AddUnlinkedNode(unsafe.Pointer(nodePtr))
		}
	*/
	return
}

func (xmlNode *XmlNode) addNextSibling(node Node) (err error) {
	nodeType := node.NodeType()
	if nodeType == XML_DOCUMENT_NODE || nodeType == XML_HTML_DOCUMENT_NODE {
		err = ERR_CANNOT_MAKE_DUCMENT_AS_CHILD
		return
	}
	nodePtr := node.NodePtr()
	if xmlNode.NodePtr() == nodePtr {
		return
	}
	ret := xmlNode.isAccestor(nodePtr)
	if ret < 0 {
		return
	} else if ret == 0 {
		if !xmlNode.Document.RemoveUnlinkedNode(nodePtr) {
			C.xmlUnlinkNodeWithCheck((*C.xmlNode)(nodePtr))
		}
		C.xmlAddNextSibling(xmlNode.Ptr, (*C.xmlNode)(nodePtr))
	} else if ret > 0 {
		node.Remove()
	}
	/*
		childPtr := C.xmlAddNextSibling(xmlNode.Ptr, (*C.xmlNode)(nodePtr))
		if nodeType == XML_TEXT_NODE && childPtr != (*C.xmlNode)(nodePtr) {
			//check the retured pointer
			//if it is not the text node just added, it means that the text node is freed because it has merged into other nodes
			//then we should invalid this node, because we do not want to have a dangling pointer
			//node.Remove()
		}
	*/
	return
}

func (xmlNode *XmlNode) Wrap(data string) (err error) {
	newNodes, err := xmlNode.coerce(data)
	if err == nil && len(newNodes) > 0 {
		newParent := newNodes[0]
		xmlNode.addNextSibling(newParent)
		newParent.AddChild(xmlNode)
	}
	return
}

func (xmlNode *XmlNode) ParseFragment(input, url []byte, options int) (fragment *DocumentFragment, err error) {
	fragment, err = parsefragment(xmlNode.Document, xmlNode, input, url, options)
	return
}

//export xmlNodeWriteCallback
func xmlNodeWriteCallback(wbufferObj unsafe.Pointer, data unsafe.Pointer, data_len C.int) {
	wbuffer := (*WriteBuffer)(wbufferObj)
	offset := wbuffer.Offset

	if offset > len(wbuffer.Buffer) {
		panic("fatal error in xmlNodeWriteCallback")
	}

	buffer := wbuffer.Buffer[:offset]
	dataLen := int(data_len)

	if dataLen > 0 {
		if len(buffer)+dataLen > cap(buffer) {
			newBuffer := grow(buffer, dataLen)
			wbuffer.Buffer = newBuffer
		}
		destBufPtr := unsafe.Pointer(&(wbuffer.Buffer[offset]))
		C.memcpy(destBufPtr, data, C.size_t(dataLen))
		wbuffer.Offset += dataLen
	}
}

//export xmlUnlinkNodeCallback
func xmlUnlinkNodeCallback(nodePtr unsafe.Pointer, gonodePtr unsafe.Pointer) {
	xmlNode := (*XmlNode)(gonodePtr)
	xmlNode.Document.AddUnlinkedNode(nodePtr)
}

func grow(buffer []byte, n int) (newBuffer []byte) {
	newBuffer = makeSlice(2*cap(buffer) + n)
	copy(newBuffer, buffer)
	return
}

func makeSlice(n int) []byte {
	// If the make fails, give a known error.
	defer func() {
		if recover() != nil {
			panic(ErrTooLarge)
		}
	}()
	return make([]byte, n)
}

func (xmlNode *XmlNode) isAccestor(nodePtr unsafe.Pointer) int {
	parentPtr := xmlNode.Ptr.parent

	if C.xmlNodePtrCheck(unsafe.Pointer(parentPtr)) == C.int(0) {
		return -1
	}
	for ; parentPtr != nil; parentPtr = parentPtr.parent {
		if C.xmlNodePtrCheck(unsafe.Pointer(parentPtr)) == C.int(0) {
			return -1
		}
		p := unsafe.Pointer(parentPtr)
		if p == nodePtr {
			return 1
		}
	}
	return 0
}

func (xmlNode *XmlNode) RecursivelyRemoveNamespaces() (err error) {
	nodePtr := xmlNode.Ptr
	C.xmlSetNs(nodePtr, nil)

	for child := xmlNode.FirstChild(); child != nil; {
		child.RecursivelyRemoveNamespaces()
		child = child.NextSibling()
	}

	nodeType := xmlNode.NodeType()

	if ((nodeType == XML_ELEMENT_NODE) ||
		(nodeType == XML_XINCLUDE_START) ||
		(nodeType == XML_XINCLUDE_END)) &&
		(nodePtr.nsDef != nil) {
		C.xmlFreeNsList((*C.xmlNs)(nodePtr.nsDef))
		nodePtr.nsDef = nil
	}

	if nodeType == XML_ELEMENT_NODE && nodePtr.properties != nil {
		property := nodePtr.properties
		for property != nil {
			if property.ns != nil {
				property.ns = nil
			}
			property = property.next
		}
	}
	return
}

func (xmlNode *XmlNode) RemoveDefaultNamespace() {
	nodePtr := xmlNode.Ptr
	C.xmlRemoveDefaultNamespace(nodePtr)
}

func (xmlNode *XmlNode) SetNamespace(prefix, href string) {
	if xmlNode.NodeType() != XML_ELEMENT_NODE {
		return
	}

	prefixBytes := GetCString([]byte(prefix))
	prefixPtr := unsafe.Pointer(&prefixBytes[0])

	hrefBytes := GetCString([]byte(href))
	hrefPtr := unsafe.Pointer(&hrefBytes[0])

	ns := C.xmlNewNs(xmlNode.Ptr, (*C.xmlChar)(hrefPtr), (*C.xmlChar)(prefixPtr))
	C.xmlSetNs(xmlNode.Ptr, ns)
}
