package main

import (
	"encoding/csv" // <-- 新增: CSV 套件
	"fmt"
	"log"
	"os" // <-- 新增: 檔案系統套件
	"regexp"
	"sort"
	"strconv" //<-- 新增: 用於將 int 轉為 string
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
				Name:    course.Name,      // 組名就是課程名稱
				MinCode: course.Code,      // 因為是第一門，所以它的 Code 暫時就是最小的
				Courses: []Course{course}, // 將這門課放進籃子的 Courses 列表裡
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
		// 組內排序 (按 Code 的「數值」)
		sort.Slice(group.Courses, func(i, j int) bool {
			codeI, _ := strconv.Atoi(group.Courses[i].Code)
			codeJ, _ := strconv.Atoi(group.Courses[j].Code)
			return codeI < codeJ
		})
		groupsToSort = append(groupsToSort, *group)
	}

	// 步驟 3: 組間排序 (按 MinCode)
	sort.Slice(groupsToSort, func(i, j int) bool {
		minCodeI, _ := strconv.Atoi(groupsToSort[i].MinCode)
		minCodeJ, _ := strconv.Atoi(groupsToSort[j].MinCode)

		// 現在，我們比較的是兩個整數，而不是兩個字串
		return minCodeI < minCodeJ
		// return groupsToSort[i].MinCode < groupsToSort[j].MinCode
	})

	// 步驟 4: 建立最終排序好的列表
	finalSortedCourses := make([]Course, 0, len(courses))
	for _, group := range groupsToSort {
		finalSortedCourses = append(finalSortedCourses, group.Courses...)
	}

	return finalSortedCourses
}

// stripTags 是一個輔助函式，用來移除字串中所有的 HTML 標籤
func stripTags(html string) string {
	re := regexp.MustCompile(`<.*?>`)
	return re.ReplaceAllString(html, "")
}

// 回傳 error，讓呼叫者決定如何處理
func write_csv_file(fileName string, contents []Course, header []string) error {
	// 1. 檢查輸入，回傳錯誤而不是 panic
	if fileName == "" {
		return fmt.Errorf("檔名不可為空") // 使用 fmt.Errorf 建立一個新的 error
	}
	if len(contents) == 0 {
		return fmt.Errorf("內容不可為空")
	}
	if len(header) == 0 {
		return fmt.Errorf("標頭不可為空")
	}
	// 這個檢查邏輯有問題，我們應該比較 Course 結構的欄位數和 header 的長度
	// reflect.TypeOf(Course{}).NumField() vs len(header)
	// 但為了簡化，我們先假設 header 是對的

	// 2. 建立檔案
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("建立檔案 %s 失敗: %w", fileName, err) // 使用 %w 包裝原始錯誤
	}
	// 將 defer 緊跟在資源建立之後
	defer file.Close()

	// 寫入 BOM
	file.WriteString("\xEF\xBB\xBF")

	// 3. 建立 writer
	writer := csv.NewWriter(file)
	// 將 Flush 的 defer 也提前
	defer writer.Flush()

	// 4. 寫入標頭
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("寫入標頭失敗: %w", err)
	}

	// 5. 寫入內容
	for _, course := range contents {
		row := []string{
			course.Code,
			course.Type,
			course.Name,
			course.Grade,
			course.Credits,
			strconv.Itoa(course.Hours),
			course.Instructor,
			course.Notes,
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("寫入資料列失敗: %w", err)
		}
	}

	fmt.Printf("成功將資料寫入 %s！\n", fileName)
	return nil // 一切順利，回傳 nil (代表沒有錯誤)
}

// calculateHours 解析時間字串並計算總時數
func calculateHours(timeStr string) int {
	if strings.Contains(timeStr, "無資料") {
		return 0
	}
	re := regexp.MustCompile(`\[.*?\]`) // 移除地點資訊 [H307]
	cleanedStr := re.ReplaceAllString(timeStr, "")
	re = regexp.MustCompile(`\d+`) // 找出所有數字
	matches := re.FindAllString(cleanedStr, -1)
	return len(matches)
}
func ParseGradeFromNotes(notes string) string {
	// 使用正則表達式尋找匹配項
	matches := gradeRegex.FindStringSubmatch(notes)

	// 如果找到了匹配項 (長度為 2：全匹配 + 1個捕獲組)
	if len(matches) == 2 {
		// matches[1] 就是我們捕獲的班級資訊，例如 "1A" 或 "3,4"
		// 清理一下頭尾可能的多餘空白後回傳
		return strings.TrimSpace(matches[1])
	}else{
		log.Printf("other info")
		for _,content := range matches{
			log.Println(content, notes)
		}
	}

	// 如果找不到，回傳一個預設值 "未知"
	return "未知"
}
