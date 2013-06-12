package utils

import (
	"strings"
)

var (
	searchengineBot []string = []string{
		"baiduspider",
		"googlebot",
		"mediapartners-google",
		"msnbot",
		"yodaobot",
		"sosospider+",
		// "yahoo! slurp;",
		// "yahoo! slurp china;",
		"yahoo! slurp",
		"iaskspider",
		"sogou spider",
		"sogou web spider",
		"sogou push spider",
	}
)

// 检查是否为搜索引擎爬虫
func IsSpider(userAgent string) bool {
	userAgent = strings.ToLower(userAgent)
	for _, v := range searchengineBot {
		if strings.Contains(userAgent, v) {
			return true
		}
	}
	return false
}
