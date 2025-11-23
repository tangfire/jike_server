package utils

import "strings"

// JoinToString 将字符串数组拼接成用逗号分隔的字符串
func JoinToString(arr []string) string {
	return strings.Join(arr, ",")
}

// ParseToArray 将逗号分隔的字符串解析为字符串数组
func ParseToArray(str string) []string {
	if str == "" {
		return []string{}
	}
	return strings.Split(str, ",")
}
