package xml

import "testing"
import "fmt"

func TestAddChild(t *testing.T) {

	docAssertion := func(doc *XmlDocument) (string, string, string) {
		expectedDocAfterAdd :=
			`<?xml version="1.0" encoding="utf-8"?>
<foo>
  <bar/>
</foo>
`
		doc.Root().AddChild("<bar></bar>")

		return doc.String(), expectedDocAfterAdd, "output of the xml doc after AddChild does not match"
	}

	nodeAssertion := func(doc *XmlDocument) (string, string, string) {
		expectedNodeAfterAdd :=
			`<foo>
  <bar/>
</foo>`

		return doc.Root().String(), expectedNodeAfterAdd, "the output of the xml root after AddChild does not match"
	}

	RunTest(t, "node", "add_child", nil, docAssertion, nodeAssertion)

}

func TestAddAncestorAsChild(t *testing.T) {
	docAssertion := func(doc *XmlDocument) (string, string, string) {
		expectedDocAfterAdd :=
			`<?xml version="1.0" encoding="utf-8"?>
<foo/>
`

		foo := doc.Root()
		bar := foo.FirstChild()
		holiday := bar.FirstChild()
		fun := holiday.FirstChild()
		fun.AddChild(bar)

		return doc.String(), expectedDocAfterAdd, "output of the xml doc after AddChild does not match"
	}

	nodeAssertion := func(doc *XmlDocument) (string, string, string) {
		expectedNodeAfterAdd :=
			`<foo/>`

		return doc.Root().String(), expectedNodeAfterAdd, "the output of the xml root after AddChild does not match"
	}

	RunTest(t, "node", "add_ancestor", nil, docAssertion, nodeAssertion)

}

func addChildBenchLogic(b *testing.B, doc *XmlDocument) {
	root := doc.Root()

	for i := 0; i < b.N; i++ {
		root.AddChild("<bar></bar>")
	}
}

func BenchmarkAddChild(b *testing.B) {
	RunBenchmark(b, "document", "big_un", addChildBenchLogic) // Run against big doc
}

func BenchmarkAddChildBigDoc(b *testing.B) {
	RunBenchmark(b, "node", "add_child", addChildBenchLogic)
}

func TestAddPreviousSibling(t *testing.T) {

	testLogic := func(t *testing.T, doc *XmlDocument) {
		err := doc.Root().AddPreviousSibling("<bar></bar><cat></cat>")

		if err != nil {
			t.Errorf("Error adding previous sibling:\n%v\n", err.Error())
		}
	}

	RunTest(t, "node", "add_previous_sibling", testLogic)
}

func TestAddPreviousSibling2(t *testing.T) {

	testLogic := func(t *testing.T, doc *XmlDocument) {
		err := doc.Root().FirstChild().AddPreviousSibling("COOL")

		if err != nil {
			t.Errorf("Error adding previous sibling:\n%v\n", err.Error())
		}
	}

	RunTest(t, "node", "add_previous_sibling2", testLogic)
}

func TestAddNextSibling(t *testing.T) {

	testLogic := func(t *testing.T, doc *XmlDocument) {
		doc.Root().AddNextSibling("<bar></bar><baz></baz>")
	}

	RunTest(t, "node", "add_next_sibling", testLogic)
}

func TestSetContent(t *testing.T) {

	testLogic := func(t *testing.T, doc *XmlDocument) {
		root := doc.Root()
		root.SetContent("<fun></fun>")
	}

	RunTest(t, "node", "set_content", testLogic)
}

func BenchmarkSetContent(b *testing.B) {

	benchmarkLogic := func(b *testing.B, doc *XmlDocument) {
		root := doc.Root()
		for i := 0; i < b.N; i++ {
			root.SetContent("<fun></fun>")
		}
	}

	RunBenchmark(b, "node", "set_content", benchmarkLogic)
}

func TestSetChildren(t *testing.T) {
	testLogic := func(t *testing.T, doc *XmlDocument) {
		root := doc.Root()
		root.SetChildren("<fun></fun>")
	}

	RunTest(t, "node", "set_children", testLogic)
}

func BenchmarkSetChildren(b *testing.B) {
	benchmarkLogic := func(b *testing.B, doc *XmlDocument) {
		root := doc.Root()
		for i := 0; i < b.N; i++ {
			root.SetChildren("<fun></fun>")
		}
	}

	RunBenchmark(b, "node", "set_children", benchmarkLogic)
}

func TestReplace(t *testing.T) {

	testLogic := func(t *testing.T, doc *XmlDocument) {
		root := doc.Root()
		root.Replace("<fun></fun><cool/>")
	}

	rootAssertion := func(doc *XmlDocument) (string, string, string) {
		root := doc.Root()
		return root.String(), "<fun/>", "the output of the xml root does not match"
	}

	RunTest(t, "node", "replace", testLogic, rootAssertion)
}

func BenchmarkReplace(b *testing.B) {

	benchmarkLogic := func(b *testing.B, doc *XmlDocument) {
		root := doc.Root()
		for i := 0; i < b.N; i++ {
			root.Replace("<fun></fun>")
			root = doc.Root() //once the node has been replaced, we need to get a new node
		}
	}

	RunBenchmark(b, "node", "replace", benchmarkLogic)
}

func TestAttributes(t *testing.T) {

	testLogic := func(t *testing.T, doc *XmlDocument) {

		root := doc.Root()
		attributes := root.Attributes()

		if len(attributes) != 2 || attributes["myname"].String() != "ff" {
			fmt.Printf("%v, %q\n", attributes, attributes["myname"].String())
			t.Error("root's attributes do not match")
		}

		child := root.FirstChild()
		childAttributes := child.Attributes()

		if len(childAttributes) != 1 || childAttributes["class"].String() != "shine" {
			t.Error("child's attributes do not match")
		}
	}

	RunTest(t, "node", "attributes", testLogic)

}

func BenchmarkAttributes(b *testing.B) {
	benchmarkLogic := func(b *testing.B, doc *XmlDocument) {

		root := doc.Root()

		for i := 0; i < b.N; i++ {
			root.SetAttr("garfield", "spaghetti")
		}
	}

	RunBenchmark(b, "node", "attributes", benchmarkLogic)
}

func TestInner(t *testing.T) {

	testLogic := func(t *testing.T, doc *XmlDocument) {
		root := doc.Root()
		root.SetInnerHtml("<bar></bar><baz></baz>")
	}

	RunTest(t, "node", "inner", testLogic)
}
func TestInnerWithAttributes(t *testing.T) {

	testLogic := func(t *testing.T, doc *XmlDocument) {
		root := doc.Root()
		root.SetInnerHtml("<bar give='me' something='good' to='eat'></bar>")
	}

	RunTest(t, "node", "inner_with_attributes", testLogic)
}

func TestSetNamespace(t *testing.T) {
	testLogic := func(t *testing.T, doc *XmlDocument) {
		root := doc.Root()
		root.SetNamespace("foo", "bar")
	}

	RunTest(t, "node", "set_namespace", testLogic)
}
