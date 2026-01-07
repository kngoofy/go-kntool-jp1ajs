package main

import (
	"fmt"
	"math"
	"unicode/utf8"

	"github.com/xuri/excelize/v2"
)

// 列幅の自動調整を行うかどうかのフラグ
// const isAutoFitColumn bool = false
var isAutoFitColumn bool = false

// CoordsToCellRef: 列インデックスと行インデックスからセルアドレスに変換
// colIndex: 列番号(1から始まる)
// rowIndex: 行番号(1から始まる)
// 戻り値: A1形式のセルアドレス文字列(例: "A1", "B5"等)
func CoordsToCellRef(colIndex, rowIndex int) string {
	col, _ := excelize.ColumnNumberToName(colIndex)
	return fmt.Sprintf("%s%d", col, rowIndex)
}

// autoFitColumn: 指定列の幅を内容に合わせて自動調整
// f: Excelファイルオブジェクト
// sheet: シート名
// col: 列アドレス(例: "A")
// 戻り値: エラー(あれば非nil以外の値)
// func autoFitColumn(f *excelize.File, sheet string, col string) error {
// 	rows, err := f.GetRows(sheet)
// 	if err != nil {
// 		return err
// 	}

// 	maxLen := 0
// 	colIndex, err := excelize.ColumnNameToNumber(col)
// 	if err != nil {
// 		return err
// 	}

// 	// 指定列の最大文字数(日本語対応)を計測
// 	for _, row := range rows {
// 		if colIndex-1 < len(row) {
// 			cellValue := row[colIndex-1]
// 			length := utf8.RuneCountInString(cellValue)
// 			if length > maxLen {
// 				maxLen = length
// 			}
// 		}
// 	}

//		// Excelの列幅は文字数に比例するため、
//		// 最大文字数に2を加えて余裕を持たせ、100を上限とする
//		width := math.Min(float64(maxLen)+2, 100)
//		return f.SetColWidth(sheet, col, col, width)
//	}
func autoFitColumnFast(f *excelize.File, sheet string, cols []string) error {

	maxLens := make([]int, len(cols))

	rows, err := f.GetRows(sheet)
	if err != nil {
		return err
	}

	// maxLen := 0
	// colIndex, err := excelize.ColumnNameToNumber(col)
	// if err != nil {
	// 	return err
	// }

	// 指定列の最大文字数(日本語対応)を計測
	for _, row := range rows {
		for colno, col := range cols {

			colIndex, err := excelize.ColumnNameToNumber(col)
			if err != nil {
				return err
			}

			if colIndex-1 < len(row) {
				cellValue := row[colIndex-1]

				length := utf8.RuneCountInString(cellValue)
				// length := len(////////////////////////////////////////cellValue)

				// if length > maxLen {
				// 	maxLen = length
				// }

				if length > maxLens[colno] {
					maxLens[colno] = length
				}
			}
		}
	}

	for colno, col := range cols {
		width := math.Min(float64(maxLens[colno])+3, 100)
		f.SetColWidth(sheet, col, col, width)
	}

	return nil
}

func setStyle(f *excelize.File) map[string]int {

	styles := map[string]int{}

	style1, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "0000FF", Style: 1},
			{Type: "top", Color: "0000FF", Style: 1},
			{Type: "bottom", Color: "0000FF", Style: 1},
			{Type: "right", Color: "0000FF", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Family: "Meiryo UI",
			Size:   11,
		},
	})
	if err != nil {
		fmt.Println(err)
	}

	style2, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "0000FF", Style: 1},
			{Type: "top", Color: "0000FF", Style: 1},
			{Type: "bottom", Color: "0000FF", Style: 1},
			{Type: "right", Color: "0000FF", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Family: "BIZ UDゴシック",
			Size:   12,
		},
	})
	if err != nil {
		fmt.Println(err)
	}

	styles["style1"] = style1
	styles["style2"] = style2

	return styles

}
