package xml

type AttributeNode struct {
	*XmlNode
}

func (attrNode *AttributeNode) String() string {
	return attrNode.Content()
}

func (attrNode *AttributeNode) Value() string {
	return attrNode.Content()
}

func (attrNode *AttributeNode) SetValue(val interface{}) {
	attrNode.SetContent(val)
}

/*
alias :value :content
alias :to_s :content
alias :content= :value=
*/
