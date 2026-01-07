package main

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

// update_sheet_def: 詳細データ(DEF)をExcelシートに書き込む共通関数
// この関数は全てのDEFシート(DefSnd, DefRcv, DefJob等)で使用される汎用関数です。
//
// 処理フロー:
// 1. dataの各要素をシートの4行目以降に書き込む
// 2. B列(連番)にstyle1を適用
// 3. C列(詳細情報)にstyle2を適用
//
// パラメータ:
//
//	f: Excelファイルオブジェクト
//	data: 詳細データの配列(連番と詳細情報のペア)
//	sheet: 対象シート名
//	styles: スタイルIDを格納するマップ
//	  - "style1": 連番列(B列)用スタイルID
//	  - "style2": 詳細情報列(C列)用スタイルID
func update_sheet_ajsprint(f *excelize.File, data []string, styles map[string]int) {

	const sheet = "Ajsprint"
	var cellref string = ""

	// 各詳細データレコードをシートの4行目以降に書き込む
	// ループ: i: 配列のインデックス(0から開始), v: 詳細データの内容
	for i, v := range data {
		// 書き込み開始セルを計算(2列目から始まる)
		// 例: i=0の場合は B4, i=1の場合は B5 となる
		cellref = CoordsToCellRef(2, i+4)

		// 詳細データレコード: 連番と詳細情報のペアを作成
		rec := []interface{}{
			i + 1, // 連番 (1から始まる)
			v,     // 詳細データ内容
		}

		// Excelシートに行データを書き込む
		// エラーが発生した場合は処理を中断して終了
		err := f.SetSheetRow(sheet, cellref, &rec)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// B列(連番列)のデータ領域にスタイル1(style1)を適用
	// 計算: 最後の行 = レコード数 + 4(ヘッダー行) - 1
	// 例: レコード数が10の場合、最後の行は13となる(4+10-1=13)
	endcell := fmt.Sprintf("B%d", len(data)+4-1)
	// fmt.Println(endcell) // デバッグ用コメントアウト
	err := f.SetCellStyle(sheet, "B4", endcell, styles["style1"])
	if err != nil {
		fmt.Println(err)
	}

	// C列(詳細データ列)のデータ領域にスタイル2(style2)を適用
	// B列と同じ範囲(行)に対してスタイルを適用
	endcell = fmt.Sprintf("C%d", len(data)+4-1)
	// fmt.Println(endcell) // デバッグ用コメントアウト
	err = f.SetCellStyle(sheet, "C4", endcell, styles["style2"])
	if err != nil {
		fmt.Println(err)
	}
}
