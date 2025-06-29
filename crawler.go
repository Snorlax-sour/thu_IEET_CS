package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"strconv"

	"github.com/PuerkitoBio/goquery" // <--- 1. 新增 goquery 套件
	"github.com/gocolly/colly/v2"
)

var gradeRegex = regexp.MustCompile(`\/\s*資工系\s*([A-Z0-9,]+)`)

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

func main() {
	var year, term string
	year_int := 108

	for year_int <= 113 {
		term_int := 1
		for term_int <= 2 {

			year = strconv.Itoa(year_int)
			term = strconv.Itoa(term_int)
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
				course.Code = strings.Join(strings.Fields(course.Code), " ") // field這是去除中間空格等部份，然後join將結果拼接起來
				// fmt.Println("code : %s", course.Code)
				// for _, code := range course.Code {
				// 	fmt.Println(code)
				// }

				log.Printf("course.code = %s, should not have space", course.Code)

				// 優化課程名稱處理：將 <br> 換成空格
				courseNameHTML, _ := e.DOM.Find("td[data-title='課程名稱'] > a").Html()
				courseNameHTML = strings.TrimSpace(strings.ReplaceAll(courseNameHTML, "<br/>", " "))
				course.Type, course.Name = Split_course_type_and_name(courseNameHTML)

				course.Credits = strings.TrimSpace(e.DOM.Find("td[data-title='學分數']").Text())

				var instructors []string
				// --- 2. 修正此處的類型 ---
				e.DOM.Find("td[data-title='授課教師'] a").Each(func(i int, s *goquery.Selection) {
					instructors = append(instructors, strings.TrimSpace(s.Text()))
				})
				// for _, content := range instructors {
				// fmt.Println(len(instructors))
				// fmt.Printf("instuctor %s, ", content)
				// }
				// log.Printf("\n")

				
				if len(instructors) > 1 && instructors[1] != "" { // i.e. above 2 people
					// need except '' symbol

					course.Instructor = strings.Join(instructors, "與")
					// for _, content := range instructors {
					// 	fmt.Printf("instructor = '%s', ", content)
					// 	for __, char := range content {
					// 		fmt.Printf("%d, %c ", __, char)

					// 	}
					// }
					// fmt.Println()
				} else {
					// log.Printf("year = %s, term = %s, len instructors = %d\n", year, term, len(instructors))
					// for _, content := range instructors {
					// log.Println(content)
					// }
					// log.Println()
					course.Instructor = instructors[0]
				}
				course.Instructor += "教授"
				log.Printf("year = %s, term = %s, instuctor = %s", year, term, course.Instructor)
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
			if term == "1" {
				term = "上"
			} else if term == "2" {
				term = "下"
			}
			// --- 在爬取完成後，寫入檔案前，呼叫新的排序函式 ---
			fmt.Printf("\n現正進行排序...\n")
			courses = SortCoursesAsTeam(courses) // << --- 在這裡呼叫！
			fmt.Println("排序完成！")

			file_name := year + "學年" + term + "課程紀錄" + ".csv"
			write_csv_file(file_name, courses, header)
			term_int++ // forgot add
		}
		year_int++ // forgot

	}
	// fmt.Print("請輸入學年 (例如: 112): ")
	// if _, err := fmt.Scanln(&year); err != nil {
	// 	log.Fatalf("讀取學年時發生錯誤: %v", err)
	// }

	// fmt.Print("請輸入學期 (1 為上學期, 2 為下學期): ")
	// if _, err := fmt.Scanln(&term); err != nil {
	// 	log.Fatalf("讀取學期時發生錯誤: %v", err)
	// }

	// if term != "1" && term != "2" {
	// 	log.Fatal("學期輸入無效，請輸入 1 或 2。")
	// }

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
