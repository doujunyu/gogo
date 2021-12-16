package utility

import "strconv"

// StringBySliceString 字符串类型切片转成字符串
func StringBySliceString(separator string, slice []string) string {
	if len(slice) == 0 {
		return ""
	}
	set := ""
	for _, val := range slice {
		set += val + separator
	}
	return set[:len(set)-len(separator)]
}

// StringBySliceInt 数字切片转成字符串
func StringBySliceInt(separator string, slice []int) string {
	set := ""
	for _, val := range slice {
		set += strconv.Itoa(val) + separator
	}
	return set[:len(set)-len(separator)]
}
