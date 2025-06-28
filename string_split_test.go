package main

import (
	"reflect"
	// "sort"
	"testing"
)

func TestSplit_string(t *testing.T) {
	// 準備輸入和預期輸出
	input := "選修-這是中文課程 course name"
	expectedType := "選修"
	expectedName := "這是中文課程" // 假設您的函式會去掉英文

	// 執行函式
	actualType, actualName := Split_course_type_and_name(input)

	// 使用 t.Logf 顯示結果，方便 -v 模式下觀察
	t.Logf("輸入: %q", input)
	t.Logf("實際得到 -> 類型: %q, 名稱: %q", actualType, actualName)
	t.Logf("預期得到 -> 類型: %q, 名稱: %q", expectedType, expectedName)

	// 斷言 1: 驗證課程類型是否完全符合預期
	if actualType != expectedType {
		t.Errorf("課程類型不符：預期得到 %q, 實際得到 %q", expectedType, actualType)
	}

	// 斷言 2: 驗證課程名稱是否完全符合預期
	if actualName != expectedName {
		t.Errorf("課程名稱不符：預期得到 %q, 實際得到 %q", expectedName, actualName)
	}
}

func TestFile_writer(f *testing.T) {
	content := []Course{
		{
			Code:       "1051",
			Type:       "必修",
			Name:       "專題實作",
			Grade:      "4",
			Instructor: "呂芳懌與黃宜豊...",
			Credits:    "2-0", // Credits 是 string
			Hours:      2,     // Hours 是 int
			Notes:      "",    // Notes 是 string
		},
	}
	header := []string{"選課代碼", "課程類別", "課程名稱", "年級/班級", "學分數", "時數", "授課教師", "備註"}

	// --- 測試案例的邏輯也需要修正 ---
	test1 := write_csv_file("", content, header)
	if test1 == nil {
		f.Errorf("這是一個 Bug！我預期它失敗，它卻成功了！") // 「那麼，我的測試就失敗了！」
	}
	/*
		一個好的心法：在寫每一個測試斷言 (assertion) 時，先問自己一句話：
		「在這個測試案例下，我預期函式回傳的 err 應該是什麼？ nil 還是 non-nil？」

	*/
	// if test1 == true { ... }

	// 測試 2: 寫入多行資料
	// 您的 have_two_row_test 也有同樣的欄位順序/類型錯誤，我們用具名欄位修正
	have_two_row_test := []Course{
		{Code: "1051", Type: "必修", Name: "專題實-1", Instructor: "教授群A", Credits: "2-0", Hours: 2, Notes: "備註A"},
		{Code: "1052", Type: "選修", Name: "專題實-2", Instructor: "教授群B", Credits: "3-0", Hours: 3, Notes: "備註B"},
	}

	// 您的測試邏輯反了，應該是預期 write_csv_file 回傳 true 才對
	success := write_csv_file("test.csv", have_two_row_test, header)
	if success != nil {
		f.Errorf(" write_csv_file() 能夠寫入多行資料時遇到問題 %s", success)
	}
	// 這裡可以加上檢查檔案是否真的被建立、內容是否正確等更進階的測試
}

func TestParseGradeFromNotes(t *testing.T) {
	// 1. 定義一個結構體來存放我們的測試案例，這樣更清晰
	testCases := []struct {
		name     string // 測試案例的名稱
		input    string // 輸入的備註字串 (notes)
		expected string // 預期提取出的年級/班級
	}{
		// --- 標準情況 ---
		{
			name:     "Standard Grade 1A",
			input:    `28019 資訊工程學系  / 資工系1A<BR>`,
			expected: "1A",
		},
		{
			name:     "Multiple Grades with Comma",
			input:    `28240 資訊工程學系  / 資工系2A,2B<BR>2A、2B併班`,
			expected: "2A,2B",
		},
		{
			name:     "Multiple Grades with Space",
			input:    `28022 資訊工程學系  / 資工系3,4<BR>限重修生`,
			expected: "3,4",
		},
		{
			name:     "Single Digit Grade",
			input:    `22107 資訊工程學系  / 資工系3<BR>人工選課`,
			expected: "3",
		},
		// --- 邊界/特殊情況 ---
		{

			name:     "Master's Program Grade Range should be ignored",
			input:    `68001 資訊工程學系  / 資工碩2-4<BR>`,
			expected: "未知", // <--- 預期它被忽略，回傳 "未知"

		},
		{
			name:     "No space after slash",
			input:    `28019 資訊工程學系/資工系1A<BR>`,
			expected: "1A",
		},
		{
			name:     "No trailing HTML tag",
			input:    `28019 資訊工程學系  / 資工系1A`,
			expected: "1A",
		},
		{
			name:     "No match found",
			input:    `這是一個沒有班級資訊的備註`,
			expected: "未知", // 預期回傳預設值
		},
		{
			name:     "Empty input string",
			input:    "",
			expected: "未知",
		},
	}

	// 2. 遍歷所有測試案例
	for _, tc := range testCases {
		// 使用 t.Run 可以讓每個測試案例在報告中獨立顯示，非常方便
		t.Run(tc.name, func(t *testing.T) {
			// 執行待測函式
			actual := ParseGradeFromNotes(tc.input)

			// 驗證實際結果是否與預期相符
			if actual != tc.expected {
				t.Errorf("對於輸入 %q:\n預期得到 %q, \n實際得到 %q", tc.input, tc.expected, actual)
			}
		})
	}
}
func TestSortCourses(t *testing.T) {
	// 1. 準備亂序的測試資料
	unsortedCourses := []Course{
		{Name: "作業系統", Code: "1039"},
		{Name: "演算法", Code: "1001"},
		{Name: "作業系統", Code: "1040"},
		{Name: "作業系統", Code: "1037"},
		{Name: "編譯器", Code: "1043"},
	}

	// 2. 定義符合「正確中文排序」和「代碼排序」的預期順序
	expectedOrder := []Course{
		{Name: "演算法", Code: "1001"},
		{Name: "作業系統", Code: "1037"},
		{Name: "作業系統", Code: "1039"},
		{Name: "作業系統", Code: "1040"},
		{Name: "編譯器", Code: "1043"},
	}

	// 3. 呼叫我們獨立出來的 SortCourses 函式

	unsortedCourses = SortCoursesAsTeam(unsortedCourses) // << --- 呼叫被測試的函式

	// 4. 驗證結果
	if !reflect.DeepEqual(unsortedCourses, expectedOrder) {
		t.Errorf("排序結果不符合預期！")
		t.Logf("--- 預期順序 ---")
		for _, c := range expectedOrder {
			t.Logf("%+v", c)
		}
		t.Logf("--- 實際順序 ---")
		for _, c := range unsortedCourses {
			t.Logf("%+v", c)
		}
	}
}
