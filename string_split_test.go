package main

import (
	"os"
	"reflect"
	"testing"
)

// 測試 1: Split_course_type_and_name
func TestSplit_course_type_and_name(t *testing.T) {
	testCases := []struct {
		name         string
		input        string
		expectedType string
		expectedName string
	}{
		{
			name:         "Standard case with English",
			input:        "選修-這是中文課程 course name",
			expectedType: "選修",
			expectedName: "這是中文課程",
		},
		{
			name:         "No English part",
			input:        "必修-資料結構",
			expectedType: "必修",
			expectedName: "資料結構",
		},
		{
			name:         "No separator",
			input:        "專題討論",
			expectedType: "未知",
			expectedName: "專題討論",
		},
		{
			name:         "Empty string",
			input:        "",
			expectedType: "未知",
			expectedName: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualType, actualName := Split_course_type_and_name(tc.input)
			if actualType != tc.expectedType {
				t.Errorf("Type mismatch: expected %q, got %q", tc.expectedType, actualType)
			}
			if actualName != tc.expectedName {
				t.Errorf("Name mismatch: expected %q, got %q", tc.expectedName, actualName)
			}
		})
	}
}

// 測試 2: SortCoursesAsTeam (自行車隊排序法)
func TestSortCoursesAsTeam(t *testing.T) {
	// 準備亂序的測試資料
	unsortedCourses := []Course{
		{Name: "B", Code: "30"},
		{Name: "A", Code: "20"},
		{Name: "C", Code: "5"},
		{Name: "B", Code: "10"},
	}

	// 根據「組的最小代碼」進行手動排序，得出預期結果
	// A組: MinCode=20
	// B組: MinCode=10
	// C組: MinCode=5
	// 組的排序應該是 C -> B -> A
	expectedOrder := []Course{
		{Name: "C", Code: "5"},  // C組
		{Name: "B", Code: "10"}, // B組 (內部按Code排序)
		{Name: "B", Code: "30"},
		{Name: "A", Code: "20"}, // A組
	}

	// 呼叫被測試的函式
	actualSorted := SortCoursesAsTeam(unsortedCourses)

	// 驗證結果
	if !reflect.DeepEqual(actualSorted, expectedOrder) {
		t.Errorf("SortCoursesAsTeam() result does not match expected order!")
		t.Logf("--- EXPECTED ORDER ---")
		for _, c := range expectedOrder {
			t.Logf("%+v", c)
		}
		t.Logf("--- ACTUAL ORDER ---")
		for _, c := range actualSorted {
			t.Logf("%+v", c)
		}
	}
}

// 測試 3: write_csv_file
func Test_write_csv_file(t *testing.T) {
	header := []string{"選課代碼", "課程類別", "課程名稱", "年級/班級", "學分數", "時數", "授課教師", "備註"}
	content := []Course{
		{Code: "1001", Type: "必修", Name: "測試課", Grade: "1", Credits: "3-0", Hours: 3, Instructor: "測試員", Notes: "無"},
	}

	// 子測試 1: 成功寫入的案例
	t.Run("successful write", func(t *testing.T) {
		fileName := "test_success.csv"

		// t.Cleanup() 是一個非常有用的函式，它會在該測試（或子測試）結束後自動執行
		// 我們用它來確保測試檔案一定會被刪除，無論測試是成功還是失敗
		t.Cleanup(func() {
			os.Remove(fileName)
		})

		err := write_csv_file(fileName, content, header)
		if err != nil {
			t.Errorf("Expected to write file successfully, but got error: %v", err)
		}

		// (可以加入更多檢查，例如讀取檔案內容進行比對)
	})

	// 子測試 2: 檔名為空的失敗案例
	t.Run("fail with empty filename", func(t *testing.T) {
		err := write_csv_file("", content, header)
		if err == nil {
			t.Error("Expected an error for empty filename, but got nil")
		}
	})

	// 子測試 3: 內容為空的失敗案例
	t.Run("fail with empty content", func(t *testing.T) {
		fileName := "test_fail_empty_content.csv"
		t.Cleanup(func() {
			os.Remove(fileName)
		})

		err := write_csv_file(fileName, []Course{}, header)
		if err == nil {
			t.Error("Expected an error for empty content, but got nil")
		}
	})
}

// 測試 4: stripTags
func TestStripTags(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple tag", "<p>Hello</p>", "Hello"},
		{"nested tags", "<div><span>World</span></div>", "World"},
		{"string with no tags", "Just a string", "Just a string"},
		{"empty string", "", ""},
		{"tag with attributes", `<a href="/path">Link</a>`, "Link"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if actual := stripTags(tc.input); actual != tc.expected {
				t.Errorf("stripTags(%q): expected %q, got %q", tc.input, tc.expected, actual)
			}
		})
	}
}

// 測試 5: calculateHours
func TestCalculateHours(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected int
	}{
		{"standard three hours", "星期二/6,7,8[C118]", 3},
		{"multiple days", "星期二/9,三/9[H307] 星期一/1,2[H308]", 4},
		{"no location info", "一/1,2,3", 3},
		{"no data string", "無資料", 0},
		{"empty string", "", 0},
		{"only location", "[ST436]", 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if actual := calculateHours(tc.input); actual != tc.expected {
				t.Errorf("calculateHours(%q): expected %d, got %d", tc.input, tc.expected, actual)
			}
		})
	}
}

// 測試 6: ParseGradeFromNotes (您已提供，非常完整，直接保留)
func TestParseGradeFromNotes(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"Standard Grade 1A", `28019 資訊工程學系  / 資工系1A<BR>`, "1A"},
		{"Multiple Grades with Comma", `28240 資訊工程學系  / 資工系2A,2B<BR>2A、2B併班`, "2A,2B"},
		{"Master's Program ignored", `68001 資訊工程學系  / 資工碩2-4<BR>`, "未知"},
		{"No match found", `這是一個沒有班級資訊的備註`, "未知"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if actual := ParseGradeFromNotes(tc.input); actual != tc.expected {
				t.Errorf("ParseGradeFromNotes(%q): expected %q, got %q", tc.input, tc.expected, actual)
			}
		})
	}
}
