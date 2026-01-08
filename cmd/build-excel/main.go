// =============================================================================
// パッケージ: main
// 概要: JP1/AJSユニット定義ファイルを解析し、Excelファイルに出力するツール
// 機能:
//   - YAML設定ファイルの読み込み
//   - CLIフラグによる設定の上書き
//   - JP1/AJSユニット定義のパース
//   - Excelテンプレートへのデータ書き込み
//   - 最終Excelファイルの生成
//
// =============================================================================
package main

import (
	"flag" // コマンドラインフラグ処理用
	"fmt"  // 書式付き出力用
	"os"   // ファイル操作・プログラム終了用

	"github.com/xuri/excelize/v2" // Excelファイル操作ライブラリ

	"gopkg.in/yaml.v3" // YAML設定ファイルパース用
)

// =============================================================================
// Config: アプリケーション設定を保持する構造体
// YAML設定ファイルから読み込まれる
// =============================================================================
type Config struct {
	DefDataDir       string `yaml:"defDataDir"`       // 定義ファイルが格納されているディレクトリパス
	TemplateBookName string `yaml:"templateBookName"` // Excelテンプレートファイル名
	MakedBookName    string `yaml:"makedBookName"`    // 出力先Excelファイル名
	DefJp1AjsName    string `yaml:"defJp1AjsName"`    // JP1/AJSユニット定義ファイル名
	Debug            bool   `yaml:"debug"`            // デバッグモードフラグ
	AutoFill         bool   `yaml:"autoFill"`         // 列幅自動調整フラグ
}

// =============================================================================
// loadConfig: YAML設定ファイルを読み込み、Config構造体にパースする
// =============================================================================
// 引数:
//   - path: 設定ファイルのパス
//
// 戻り値:
//   - *Config: パースされた設定情報へのポインタ
//   - error: エラー（ファイル読み込み失敗またはYAMLパース失敗時）
func loadConfig(path string) (*Config, error) {
	// 設定ファイルを読み込み
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// YAMLをConfig構造体にアンマーシャル
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// =============================================================================
// main: プログラムのエントリーポイント
// =============================================================================
// 処理概要:
//  1. CLIフラグのパース
//  2. YAML設定ファイルの読み込み
//  3. CLIフラグによる設定の上書き
//  4. JP1/AJSユニット定義のパース
//  5. Excelテンプレートを開き、データを書き込み
//  6. 最終Excelファイルを保存
func main() {

	// =========================================================================
	// CLIフラグの定義とパース
	// =========================================================================
	// -config: 設定ファイルパス（デフォルト: ./config.yaml）
	configPath := flag.String("config", "./config.yaml", "path to config file")
	// -data: 定義ファイルディレクトリ（設定ファイルを上書き）
	defDataDirFlag := flag.String("data", "", "hulft def data directory")
	// -template: Excelテンプレートファイル（設定ファイルを上書き）
	templateFlag := flag.String("template", "", "template excel book")
	// -maked: 出力先Excelファイル（設定ファイルを上書き）
	makedBooFlag := flag.String("maked", "", "maked excel book")
	// -debug: デバッグモード有効化
	debugFlag := flag.Bool("debug", false, "enable debug mode")
	// -autoFill: 列幅自動調整有効化
	autoFillFlag := flag.Bool("autoFill", false, "enable auto fill columns")

	// 全てのフラグをパース（必須）
	flag.Parse()

	// =========================================================================
	// YAML設定ファイルの読み込み
	// =========================================================================
	cfg, err := loadConfig(*configPath)
	if err != nil {
		fmt.Println("config load error:", err)
		os.Exit(1) // 設定読み込み失敗時は終了コード1で終了
	}

	// =========================================================================
	// CLIフラグによる設定の上書き
	// フラグが指定されている場合は、YAML設定より優先
	// =========================================================================
	if *templateFlag != "" {
		cfg.TemplateBookName = *templateFlag
	}
	if *defDataDirFlag != "" {
		cfg.DefDataDir = *defDataDirFlag
	}
	if *makedBooFlag != "" {
		cfg.MakedBookName = *makedBooFlag
	}
	if *debugFlag == true {
		cfg.Debug = *debugFlag
	}
	if *autoFillFlag == true {
		cfg.AutoFill = *autoFillFlag
		isAutoFitColumn = *autoFillFlag // グローバル変数も更新
	}

	// =========================================================================
	// 設定内容の確認表示
	// =========================================================================
	fmt.Printf("DefDataDir:= %s\n", cfg.DefDataDir)       // 定義ファイルディレクトリ
	fmt.Printf("Template:= %s\n", cfg.TemplateBookName)   // テンプレートファイル
	fmt.Printf("Maked:= %s\n", cfg.MakedBookName)         // 出力ファイル
	fmt.Printf("DefJp1AjsName:= %s\n", cfg.DefJp1AjsName) // JP1/AJS定義ファイル名
	fmt.Printf("Debug:= %t\n", cfg.Debug)                 // デバッグモード
	fmt.Printf("AutoFill:= %t\n", cfg.AutoFill)           // 列幅自動調整

	// =========================================================================
	// JP1/AJSユニット定義のパース
	// =========================================================================
	// ユニット定義ファイルを読み込み、抽象構文木と行データを取得
	unitBlock, lines, flckStack, arStack, netStack, jobStack := build_jp1ajs_model(cfg.DefDataDir + cfg.DefJp1AjsName)
	_ = unitBlock // 抽象構文木（将来拡張用に保持）

	// グローバル変数に列幅自動調整フラグを設定
	isAutoFitColumn = cfg.AutoFill

	// =========================================================================
	// Excelテンプレートファイルを開く
	// =========================================================================
	f, err := excelize.OpenFile(cfg.TemplateBookName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 関数終了時にExcelファイルを確実にクローズする
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// =========================================================================
	// Excelスタイルの設定
	// =========================================================================
	// ヘッダー用・データ用のスタイルを定義
	styleMap := setStyle(f)
	_ = styleMap // スタイルマップ（各シート更新関数で使用）

	// =========================================================================
	// 各シートの更新
	// =========================================================================

	// _ = arStack
	// _ = netStack
	// _ = jobStack

	update_sheet_net(f, netStack, styleMap)
	update_sheet_job(f, jobStack, styleMap)
	update_sheet_file(f, flckStack, styleMap)

	update_sheet_next(f, arStack, styleMap)

	// JP1/AJSユニット定義プリントシートを更新
	update_sheet_ajsprint(f, lines, styleMap)

	// =========================================================================
	// Indexシートの更新
	// =========================================================================
	update_sheet_index(f)

	// =========================================================================
	// 生成されたExcelファイルを保存
	// =========================================================================
	if err := f.SaveAs(cfg.MakedBookName); err != nil {
		fmt.Println(err)
	}
}
