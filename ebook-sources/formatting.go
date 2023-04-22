package ebook

import (
	"regexp"
)

var maxFileNameLength = 30

// sanitizeName 删除字符串中的特殊字符
func sanitizeName(name string) string {
	var reNameBlacklist = regexp.MustCompile(`(&|>|<|\/|:|\n|\"|\||\?|\\|\r)*`)
	RemoveSpStr := reNameBlacklist.ReplaceAllString(name, "") //去除特殊字符串
	result := []rune(RemoveSpStr)
	return string(result)
}

// SanitizeName 删除字符串中的特殊字符
func SanitizeName(name string) string {
	result := sanitizeName(name)
	return result
}
