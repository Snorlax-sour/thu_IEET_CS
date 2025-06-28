package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery" // <--- 1. 新增 goquery 套件
	"github.com/gocolly/colly/v2"
)

// Course 結構用於儲存單一課程的資訊
type Course struct {
	Code       string
	Name       string
	Type       string
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
		got := Split_course_type_and_name(course.Name)
		course.Type = got[0]
		course.Name = got[1]

		course.Credits = strings.TrimSpace(e.DOM.Find("td[data-title='學分數']").Text())

		var instructors []string
		// --- 2. 修正此處的類型 ---
		e.DOM.Find("td[data-title='授課教師'] a").Each(func(i int, s *goquery.Selection) {
			instructors = append(instructors, s.Text())
		})
		course.Instructor = strings.Join(instructors, ", ")

		// 優化備註處理：保留換行並移除所有 HTML 標籤
		notesHTML, _ := e.DOM.Find("td[data-title='備註']").Html()
		notesWithNewlines := strings.ReplaceAll(notesHTML, "<br/>", "\n")
		course.Notes = strings.TrimSpace(stripTags(notesWithNewlines))

		timeLocation := e.DOM.Find("td[data-title='時間地點']").Text()
		course.Hours = calculateHours(timeLocation)

		courses = append(courses, course)
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

	// 輸出結果
	fmt.Println("\n--- 爬取結果 ---")
	for i, course := range courses {
		fmt.Printf("\n[%d] ==============================================\n", i+1)
		fmt.Printf("選課代碼: %s\n", course.Code)
		fmt.Printf("課程屬性: %s\n", course.Type)
		fmt.Printf("課程名稱: %s\n", course.Name)
		fmt.Printf("授課教師: %s\n", course.Instructor)
		fmt.Printf("學分:     %s\n", course.Credits)
		fmt.Printf("時數:     %d 小時\n", course.Hours)
		fmt.Printf("備註:\n---\n%s\n---\n", course.Notes)
	}
}
