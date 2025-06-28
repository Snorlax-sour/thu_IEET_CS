package main

import (
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
			Instructor: "呂芳懌與黃宜豊...",
			Credits:    "2-0", // Credits 是 string
			Hours:      2,     // Hours 是 int
			Notes:      "",    // Notes 是 string
		},
	}
	header := []string{"選課代碼", "課程類別", "課程名稱", "學分數", "時數", "授課教師", "備註"}

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
