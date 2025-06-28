package main

import (
	"fmt"
	"sort"
	"strings"
	// "golang.org/x/text/collate"
	// "golang.org/x/text/language"
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

// / CourseGroup 用於輔助排序
type CourseGroup struct {
	Name    string   // 課程名稱，也作為組名
	MinCode string   // 該組的最小代碼，用於組間排序
	Courses []Course // 該組的所有課程
}

// SortCoursesAsTeam 實現您的「自行車隊」排序邏輯
func SortCoursesAsTeam(courses []Course) []Course {
	if len(courses) == 0 {
		return courses
	}

	// 步驟 1: 分組並計算每組的最小代碼
	courseGroupsMap := make(map[string]*CourseGroup)
	for _, course := range courses {
		// 這是一個 Go 的特殊語法，稱為 "comma, ok" idiom。
    // group 會接收找到的值，exists 是一個布林值，代表是否找到了。
		group, exists := courseGroupsMap[course.Name]
		if !exists {
			// 我們需要為它建立一個新的籃子 (CourseGroup)。
			courseGroupsMap[course.Name] = &CourseGroup{
				Name:    course.Name,         // 組名就是課程名稱
				MinCode: course.Code,         // 因為是第一門，所以它的 Code 暫時就是最小的
				Courses: []Course{course},    // 將這門課放進籃子的 Courses 列表裡
			}
		} else {
			group.Courses = append(group.Courses, course)
			if course.Code < group.MinCode {
				group.MinCode = course.Code
			}
		}
	}
	/*
		+-------------------------------------------------------------+
		|                                                             |
		|  +---------------------------+   +------------------------+ |
		|  | 抽屜 "A"                  |   | 抽屜 "B"                 | |
		|  +---------------------------+   +------------------------+ |
		|  |                           |   |                        | |
		|  |  +---------------------+  |   |  +-------------------+ | |
		|  |  | 檔案夾 (CourseGroup)|  |   |  |檔案夾 (CourseGroup)| | |
		|  |  |---------------------|  |   |  |-------------------| | |
		|  |  | Name: "A"           |  |   |  | Name: "B"         | | |
		|  |  | MinCode: "20"       |  |   |  | MinCode: "10"     | | |
		|  |  |---------------------|  |   |  |-------------------| | |
		|  |  | 紙張 (Courses slice):|  |   |  |紙張 (Courses slice):| | |
		|  |  |  - {A, 20}          |  |   |  | - {B, 30}         | | |
		|  |  |                     |  |   |  | - {B, 10}         | | |
		|  |  +---------------------+  |   |  +-------------------+ | |
		|  |                           |   |                        | |
		|  +---------------------------+   +------------------------+ |
		|                                                             |
		+-------------------------------------------------------------+
			
	
	*/

	// 步驟 2: 將分組從 map 移至 slice 以便排序
	groupsToSort := make([]CourseGroup, 0, len(courseGroupsMap))
	for _, group := range courseGroupsMap {
		fmt.Println("group content = ", group)
		// 組內排序 (按 Code)
		sort.Slice(group.Courses, func(i, j int) bool {
			return group.Courses[i].Code < group.Courses[j].Code
		})
		groupsToSort = append(groupsToSort, *group)
	}

	// 步驟 3: 組間排序 (按 MinCode)
	sort.Slice(groupsToSort, func(i, j int) bool {
		return groupsToSort[i].MinCode < groupsToSort[j].MinCode
	})

	// 步驟 4: 建立最終排序好的列表
	finalSortedCourses := make([]Course, 0, len(courses))
	for _, group := range groupsToSort {
		finalSortedCourses = append(finalSortedCourses, group.Courses...)
	}

	return finalSortedCourses
}
