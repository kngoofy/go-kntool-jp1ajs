package main

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

func update_sheet_index(f *excelize.File) {

	const sheet = "Index"
	var cellref string = "F3"

	date := time.Now().Format("2006.01.02")

	// Excelシートに行データを書き込む
	err := f.SetCellValue(sheet, cellref, date)
	if err != nil {
		fmt.Println(err) // エラーメッセージを表示
		return           // エラーが発生した場合は処理を中断
	}

	style, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Family: "Meiryo UI",
			Size:   11,
		},
	})

	err = f.SetCellStyle(sheet, cellref, cellref, style)
	if err != nil {
		fmt.Println(err)
	}
}
