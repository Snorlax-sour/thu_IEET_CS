package main

import (
	"strings"
)

func Split_course_type_and_name(name string) (course_type string, course_name string) {
	// this function cut '-' symbol
	cut_symbol := "-"
	if find_cut_symbol := strings.Index(name, cut_symbol); find_cut_symbol == -1 { // 寫反條件
		// 2. 如果找不到分隔符，我們需要定義一個合理的預設行為
		return "未知", name // 回傳「未知」類型和完整的原始字串
	}
	full_string := strings.Split(name, cut_symbol)

	course_type = full_string[0]
	// don't need english course name

	course_name = strings.Split(full_string[1], " ")[0] // gemini發現寫反了
	// log.Println(course.Name)
	// Go 會自動回傳已賦值的具名回傳變數 courseType 和 courseName
	return
}
