package main

import (
	"fmt"
	"log"

	"github.com/kngoofy/go-kntool-jp1ajs/pkg/model"
	"github.com/xuri/excelize/v2"
)

// update_sheet_rcv: RCV(受信)モデル情報をExcelシートに書き込む
// f: Excelファイルオブジェクト
// style: セルスタイルID
// rcvs: RCVモデルの配列
func update_sheet_net(f *excelize.File, nets []model.NetUnit, styles map[string]int) {

	const sheet = "Net"
	var cellref string = ""

	// 各RCVレコードをシートの4行目以降に書き込む
	for i, v := range nets {
		// fmt.Println(i, v)
		// 書き込み開始セルを計算(2列目から始まる)
		cellref = CoordsToCellRef(2, i+4)

		rcv := []interface{}{
			i + 1,
			v.UnitName,
			v.Cm,
			v.Ha,
			v.Sd,
			v.St,
			v.Cy,
			v.Sh,
			v.Shd,
			v.De,
			v.UnitAbsoluteName,
		}

		// シートに行データを書き込む
		err := f.SetSheetRow(sheet, cellref, &rcv)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// データの最終セルを計算してスタイルを適用
	// 計算: 最後の行 = レコード数 + 4(ヘッダー行) - 1
	endcell := fmt.Sprintf("L%d", len(nets)+4-1)
	// fmt.Println(endcell)

	// データ領域(B4からendcell)にスタイルを適用
	err := f.SetCellStyle(sheet, "B4", endcell, styles["style1"])
	if err != nil {
		fmt.Println(err)
	}

	// 必要に応じて各列の幅を内容に合わせて自動調整
	// if isAutoFitColumn {
	// 	// 自動調整対象の列を指定
	// 	// 列A(連番)は除外
	// 	cols := []string{"B", "C", "D", "E", "F", "G", "H", "J", "K",
	// 		"L", "M", "O", "P", "Q", "R", "S", "T"}
	// 	for _, col := range cols {
	// 		if err := autoFitColumn(f, sheet, col); err != nil {
	// 			log.Fatal(err)
	// 		}
	// 	}
	// }

	if isAutoFitColumn {
		cols := []string{"B", "C", "D", "E", "F", "G", "H", "J", "K",
			"L", "M", "O", "P", "Q", "R", "S", "T"}
		if err := autoFitColumnFast(f, sheet, cols); err != nil {
			log.Fatal(err)
		}
	}

}

// // update_sheet_rcv_def: RCV詳細データ(DEF)をExcelシートに書き込む
// // f: Excelファイルオブジェクト
// // style: 標準スタイルID
// // style2: データエリア用スタイルID
// // data: RCV詳細データの配列
// func update_sheet_file_def(f *excelize.File, data []string, styles map[string]int) {
// 	const sheet = "DefRcv"
// 	// 共通の詳細データ更新関数を呼び出す
// 	update_sheet_def(f, data, sheet, styles)
// }
