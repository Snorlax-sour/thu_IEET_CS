package main

import (
	"encoding/csv" // <-- 新增: CSV 套件
	"fmt"
	"log"
	"os"     // <-- 新增: 檔案系統套件
	"strconv"// <-- 新增: 用於將 int 轉為 string
	"regexp" 
	"strings"

	"github.com/PuerkitoBio/goquery" // <--- 1. 新增 goquery 套件
	"github.com/gocolly/colly/v2"
)
var gradeRegex = regexp.MustCompile(`\/\s*資工系\s*([A-Z0-9,]+)`)

func ParseGradeFromNotes(notes string) string {
	// 使用正則表達式尋找匹配項
	matches := gradeRegex.FindStringSubmatch(notes)

	// 如果找到了匹配項 (長度為 2：全匹配 + 1個捕獲組)
	if len(matches) == 2 {
		// matches[1] 就是我們捕獲的班級資訊，例如 "1A" 或 "3,4"
		// 清理一下頭尾可能的多餘空白後回傳
		return strings.TrimSpace(matches[1])
	}

	// 如果找不到，回傳一個預設值 "未知"
	return "未知"
}
// Course 結構用於儲存單一課程的資訊
type Course struct {
	Code       string
	Name       string
	Type       string
	Grade      string // <-- 新增此欄位
	Instructor string
	Notes      string
	Credits    string
	Hours      int
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

func main() {
	var year, term string

	fmt.Print("請輸入學年 (例如: 112): ")
	if _, err := fmt.Scanln(&year); err != nil {
		log.Fatalf("讀取學年時發生錯誤: %v", err)
	}

	fmt.Print("請輸入學期 (1 為上學期, 2 為下學期): ")
	if _, err := fmt.Scanln(&term); err != nil {
		log.Fatalf("讀取學期時發生錯誤: %v", err)
	}

	if term != "1" && term != "2" {
		log.Fatal("學期輸入無效，請輸入 1 或 2。")
	}

	url := fmt.Sprintf("https://course.thu.edu.tw/view-dept/%s/%s/350/", year, term)
	fmt.Printf("\n正在爬取目標網址: %s\n\n", url)

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	var courses []Course

	c.OnHTML("table#no-more-tables > tbody > tr", func(e *colly.HTMLElement) {
		if e.DOM.Find("td[data-title='選課代碼']").Text() == "" {
			return
		}

		course := Course{}

		course.Code = strings.TrimSpace(e.DOM.Find("td[data-title='選課代碼']").Text())

		// 優化課程名稱處理：將 <br> 換成空格
		courseNameHTML, _ := e.DOM.Find("td[data-title='課程名稱'] > a").Html()
		course.Name = strings.TrimSpace(strings.ReplaceAll(courseNameHTML, "<br/>", " "))
		course.Type, course.Name = Split_course_type_and_name(course.Name)

		course.Credits = strings.TrimSpace(e.DOM.Find("td[data-title='學分數']").Text())

		var instructors []string
		// --- 2. 修正此處的類型 ---
		e.DOM.Find("td[data-title='授課教師'] a").Each(func(i int, s *goquery.Selection) {
			instructors = append(instructors, s.Text())
		})
		course.Instructor = strings.Join(instructors, "與")
		course.Instructor += "教授"
		// 優化備註處理：保留換行並移除所有 HTML 標籤
		notesHTML, _ := e.DOM.Find("td[data-title='備註']").Html()
		// 我們從包含 HTML 標籤的 notesHTML 進行解析，因為我們的正則需要它
    	course.Grade = ParseGradeFromNotes(notesHTML) 
 		// 接著才移除 HTML 標籤，填充乾淨的 Notes 欄位
		notesWithNewlines := strings.ReplaceAll(notesHTML, "<br/>", "\n")
		course.Notes = strings.TrimSpace(stripTags(notesWithNewlines))

		timeLocation := e.DOM.Find("td[data-title='時間地點']").Text()

		if graduate_course := strings.Contains(course.Notes, "碩"); graduate_course {
			return //碩士不管，已經處理完成
		}
		course.Hours = calculateHours(timeLocation)
		find_elective_course := strings.Index(course.Type, "選修")
		if find_elective_course != -1 {
			return //  選修不添加
		} else { // else must in } follow
			courses = append(courses, course)
		}

	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Printf("\n爬取完成！共找到 %d 門課程。\n", len(courses))
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("爬取網站時發生錯誤 (狀態碼: %d): %v\n", r.StatusCode, err)
		log.Println("請檢查學年學期是否正確，或網路連線是否正常。")
	})

	if err := c.Visit(url); err != nil {
		log.Fatalf("無法訪問目標網址: %v", err)
	}
	header := []string{"選課代碼", "課程類別", "課程名稱", "學分數", "時數", "授課教師", "備註"}
	var terms string
	if term == "1"{
		terms = "上"
	}else if term == "2"{
		terms = "下"
	}
	// --- 在爬取完成後，寫入檔案前，呼叫新的排序函式 ---
    fmt.Printf("\n爬取完成！共找到 %d 門課程，現正進行排序...\n", len(courses))
    courses = SortCoursesAsTeam(courses) // << --- 在這裡呼叫！
    fmt.Println("排序完成！")

	file_name := year + "學年" + terms + "課程紀錄" + ".csv"
	write_csv_file(file_name, courses, header)

	// 輸出結果
	// fmt.Println("\n--- 爬取結果 ---")
	// for _, course := range courses {
	// 	// fmt.Printf("\n[%d] ==============================================\n", _+1)
	// 	fmt.Printf("選課代碼: %s\n", course.Code)
	// 	fmt.Printf("課程屬性: %s\n", course.Type)
	// 	fmt.Printf("課程名稱: %s\n", course.Name)
	// 	fmt.Printf("授課教師: %s\n", course.Instructor)
	// 	fmt.Printf("學分:     %s\n", course.Credits)
	// 	fmt.Printf("時數:     %d 小時\n", course.Hours)
	// 	fmt.Printf("備註:\n---\n%s\n---\n", course.Notes)
	// }
	// fmt.Println("總課程數量: ", len(courses))

}
