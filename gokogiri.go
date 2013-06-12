package gokogiri

import (
	"gokogiri/html"
	"gokogiri/xml"
)

func ParseHtml(content []byte) (doc *html.HtmlDocument, err error) {
	return html.Parse(content, html.DefaultEncodingBytes, nil, html.DefaultParseOption, html.DefaultEncodingBytes)
}

func ParseXml(content []byte) (doc *xml.XmlDocument, err error) {
	return xml.Parse(content, xml.DefaultEncodingBytes, nil, xml.DefaultParseOption, xml.DefaultEncodingBytes)
}
