// =============================================================================
// パッケージ: main
// 概要: JP1/AJSユニット定義情報をExcelファイルに出力するためのユーティリティ関数群
// 依存ライブラリ: excelize/v2 - Go言語用Excel操作ライブラリ
// 主な機能:
//   - セルアドレス変換（CoordsToCellRef）
//   - 列幅の自動調整（autoFitColumnFast）
//   - スタイル設定（setStyle）
//
// =============================================================================
package main

import (
	"fmt"          // 書式付き出力用
	"math"         // 数学関数用（Min関数）
	"unicode/utf8" // UTF-8文字列処理用（日本語文字数カウント）

	"github.com/xuri/excelize/v2" // Excelファイル操作ライブラリ
)

// =============================================================================
// グローバル設定
// =============================================================================

// isAutoFitColumn: 列幅の自動調整を行うかどうかのフラグ
// true: 自動調整を行う, false: 自動調整を行わない
// const isAutoFitColumn bool = false // 定数版（コメントアウト）
var isAutoFitColumn bool = false

// =============================================================================
// CoordsToCellRef: 列インデックスと行インデックスからセルアドレスに変換する
// =============================================================================
// 引数:
//   - colIndex: 列番号（1から始まる、1=A, 2=B, ...）
//   - rowIndex: 行番号（1から始まる）
//
// 戻り値:
//   - string: A1形式のセルアドレス文字列（例: "A1", "B5", "AA100"等）
//
// 使用例:
//   - CoordsToCellRef(1, 1) -> "A1"
//   - CoordsToCellRef(3, 10) -> "C10"
func CoordsToCellRef(colIndex, rowIndex int) string {
	// excelizeのユーティリティで列番号を列名に変換（例: 1 -> "A"）
	col, _ := excelize.ColumnNumberToName(colIndex)
	// 列名と行番号を結合してセルアドレスを生成
	return fmt.Sprintf("%s%d", col, rowIndex)
}

// =============================================================================
// autoFitColumn: 指定列の幅を内容に合わせて自動調整（単一列版 - コメントアウト）
// =============================================================================
// 引数:
//   - f: Excelファイルオブジェクト
//   - sheet: シート名
//   - col: 列アドレス（例: "A"）
//
// 戻り値:
//   - error: エラー（あれば非nil）
//
// 処理概要:
//  1. シートの全行を取得
//  2. 指定列の最大文字数（日本語対応）を計測
//  3. 列幅を最大文字数+2に設定（上限100）
// =============================================================================
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
//		// 最大文字数に2を加えて余裕を持たせ、1003092上限とする
//		width := math.Min(float64(maxLen)+2, 100)
//		return f.SetColWidth(sheet, col, col, width)
//	}
//
// =============================================================================
// autoFitColumnFast: 複数列の幅を内容に合わせて一括自動調整（高速版）
// =============================================================================
// 引数:
//   - f: Excelファイルオブジェクト
//   - sheet: シート名
//   - cols: 調整対象の列アドレスのスライス（例: []string{"A", "B", "C"}）
//
// 戻り値:
//   - error: エラー（あれば非nil）
//
// 処理概要:
//  1. シートの全行を一度だけ取得（パフォーマンス向上）
//  2. 各列の最大文字数（日本語対応: utf8.RuneCountInString）を計測
//  3. 各列の幅を最大文字数+3に設定（上限100）
//
// 注意:
//   - autoFitColumn（単一列版）と比べ、行データの取得が1回で済むため高速
//
// =============================================================================
func autoFitColumnFast(f *excelize.File, sheet string, cols []string) error {
	// 各列の最大文字数を保持するスライス
	maxLens := make([]int, len(cols))

	// シートの全行データを取得（一度だけ取得して再利用）
	rows, err := f.GetRows(sheet)
	if err != nil {
		return err
	}

	// ※ 旧コード: 単一列処理用の変数（コメントアウト）
	// maxLen := 0
	// colIndex, err := excelize.ColumnNameToNumber(col)
	// if err != nil {
	// 	return err
	// }

	// =========================================================================
	// 各行・各列を走査して最大文字数を計測
	// =========================================================================
	for _, row := range rows {
		// 指定された各列について処理
		for colno, col := range cols {
			// 列名を列番号に変換（例: "A" -> 1）
			colIndex, err := excelize.ColumnNameToNumber(col)
			if err != nil {
				return err
			}

			// 列が行の範囲内にある場合のみ処理
			if colIndex-1 < len(row) {
				// セルの値を取得
				cellValue := row[colIndex-1]

				// UTF-8文字数をカウント（日本語などのマルチバイト文字に対応）
				length := utf8.RuneCountInString(cellValue)
				// ※ バイト数でのカウント（日本語非対応 - コメントアウト）
				// length := len(cellValue)

				// ※ 旧コード: 単一列処理用（コメントアウト）
				// if length > maxLen {
				// 	maxLen = length
				// }

				// 列ごとの最大文字数を更新
				if length > maxLens[colno] {
					maxLens[colno] = length
				}
			}
		}
	}

	// =========================================================================
	// 各列の幅を設定
	// =========================================================================
	for colno, col := range cols {
		// 列幅 = 最大文字数 + 3（余白）、上限100
		width := math.Min(float64(maxLens[colno])+3, 100)
		// 列幅を設定（開始列と終了列は同じ）
		f.SetColWidth(sheet, col, col, width)
	}

	return nil
}

// =============================================================================
// setStyle: Excelファイル用のスタイルを定義する
// =============================================================================
// 引数:
//   - f: Excelファイルオブジェクト
//
// 戻り値:
//   - map[string]int: スタイル名とスタイルIDのマップ
//   - "style1": 中央揃えスタイル（Meiryo UI, 11pt, 青罫線）
//   - "style2": 左揃えスタイル（BIZ UDゴシック, 12pt, 青罫線）
//
// 処理概要:
//  1. style1: ヘッダー用スタイル（中央揃え、Meiryo UIフォント）
//  2. style2: データ用スタイル（左揃え、BIZ UDゴシックフォント）
//
// =============================================================================
func setStyle(f *excelize.File) map[string]int {
	// スタイルIDを格納するマップを初期化
	styles := map[string]int{}

	// =========================================================================
	// style1: ヘッダー用スタイル
	// - 罫線: 青色（0000FF）の細線（上下左右）
	// - 配置: 水平・垂直ともに中央揃え
	// - フォント: Meiryo UI, 11pt
	// =========================================================================
	style1, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "0000FF", Style: 1},   // 左罫線
			{Type: "top", Color: "0000FF", Style: 1},    // 上罫線
			{Type: "bottom", Color: "0000FF", Style: 1}, // 下罫線
			{Type: "right", Color: "0000FF", Style: 1},  // 右罫線
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center", // 水平方向: 中央
			Vertical:   "center", // 垂直方向: 中央
		},
		Font: &excelize.Font{
			Family: "Meiryo UI", // フォントファミリー
			Size:   11,          // フォントサイズ
		},
	})
	if err != nil {
		fmt.Println(err)
	}

	// =========================================================================
	// style2: データ用スタイル
	// - 罫線: 青色（0000FF）の細線（上下左右）
	// - 配置: 水平は左揃え、垂直は中央揃え
	// - フォント: BIZ UDゴシック, 12pt
	// =========================================================================
	style2, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "0000FF", Style: 1},   // 左罫線
			{Type: "top", Color: "0000FF", Style: 1},    // 上罫線
			{Type: "bottom", Color: "0000FF", Style: 1}, // 下罫線
			{Type: "right", Color: "0000FF", Style: 1},  // 右罫線
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",   // 水平方向: 左揃え
			Vertical:   "center", // 垂直方向: 中央
		},
		Font: &excelize.Font{
			Family: "BIZ UDゴシック", // フォントファミリー（日本語フォント）
			Size:   12,           // フォントサイズ
		},
	})
	if err != nil {
		fmt.Println(err)
	}

	// スタイルをマップに登録
	styles["style1"] = style1 // ヘッダー用
	styles["style2"] = style2 // データ用

	return styles
}
