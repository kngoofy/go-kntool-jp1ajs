package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DefDataDir       string `yaml:"defDataDir"`
	TemplateBookName string `yaml:"templateBookName"`
	MakedBookName    string `yaml:"makedBookName"`
	DefJp1AjsName    string `yaml:"defJp1AjsName"`
	Debug            bool   `yaml:"debug"`
	AutoFill         bool   `yaml:"autoFill"`
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// main: プログラムのエントリーポイント
// Hulftの各種定義ファイルを読み込み、テンプレートExcelに
// データを挿入して最終的なExcelファイルを生成します。
func main() {

	// --- CLI フラグ ---
	// フラグの定義（デフォルト値・説明文・戻り値型）
	configPath := flag.String("config", "./config.yaml", "path to config file")

	defDataDirFlag := flag.String("data", "", "hulft def data directory")
	templateFlag := flag.String("template", "", "template excel book")
	makedBooFlag := flag.String("maked", "", "maked excel book")
	debugFlag := flag.Bool("debug", false, "enable debug mode")
	autoFillFlag := flag.Bool("autoFill", false, "enable auto fill columns")
	// 全てのフラグのパース（必須）
	flag.Parse()

	// --- 設定ファイル読み込み ---
	cfg, err := loadConfig(*configPath)
	if err != nil {
		fmt.Println("config load error:", err)
		os.Exit(1)
	}

	// --- フラグで上書き ---
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
		isAutoFitColumn = *autoFillFlag
	}

	fmt.Printf("DefDataDir:= %s\n", cfg.DefDataDir)
	fmt.Printf("Template:= %s\n", cfg.TemplateBookName)
	fmt.Printf("Maked:= %s\n", cfg.MakedBookName)
	fmt.Printf("DefJp1AjsName:= %s\n", cfg.DefJp1AjsName)
	fmt.Printf("Debug:= %t\n", cfg.Debug)
	fmt.Printf("AutoFill:= %t\n", cfg.AutoFill)

	unitBlock, lines := build_jp1ajs_model(cfg.DefDataDir + cfg.DefJp1AjsName)
	_ = unitBlock

	isAutoFitColumn = cfg.AutoFill

	// 各Hulft定義ファイル(SND/RCV/JOB等)からモデルとデータを読み込む
	// モデル: ファイル定義情報のモデル構造体
	// データ: ファイル内の個別データの配列
	// snd_model, rcv_model, job_model, hst_model, tgrp_model, fmt_model, mfmt_model, trg_model,
	// 	snd_data, rcv_data, job_data, hst_data, tgrp_data, fmt_data, mfmt_data, trg_data :=
	// 	build_hulft_model(
	// 		cfg.DefDataDir,
	// 		cfg.DefSndName,
	// 		cfg.DefRcvName,
	// 		cfg.DefJobName,
	// 		cfg.DefHstName,
	// 		cfg.DefTgrpName,
	// 		cfg.DefFmtName,
	// 		cfg.DefMfmtName,
	// 		cfg.DefTgrpName,
	// 	)

	// Excelテンプレートファイルを開く
	// f, err := excelize.OpenFile(templateBook)
	f, err := excelize.OpenFile(cfg.TemplateBookName)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 関数終了時にファイルを確実にクローズする
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// // スタイル設定1:
	styleMap := setStyle(f)
	_ = styleMap

	update_sheet_ajsprint(f, lines, styleMap)
	// // 各シートを更新: モデル情報をそれぞれのシートに書き込む
	// // SND(送信)シートを更新
	// update_sheet_snd(f, snd_model, styleMap)
	// // RCV(受信)シートを更新
	// update_sheet_rcv(f, rcv_model, styleMap)
	// // HST(ホスト)シートを更新
	// update_sheet_hst(f, hst_model, styleMap)
	// // TGRP(転送グループ)シートを更新
	// update_sheet_tgrp(f, tgrp_model, styleMap)
	// // JOB(ジョブ)シートを更新
	// update_sheet_job(f, job_model, styleMap)
	// // FMT(フォーマット)シートを更新
	// update_sheet_fmt(f, fmt_model, styleMap)
	// // MFMT(マッピングフォーマット)シートを更新
	// update_sheet_mfmt(f, mfmt_model, styleMap)
	// // TRG(トリガー)シートを更新
	// update_sheet_trg(f, trg_model, styleMap)

	// // 各シートの詳細データ(DEF)を更新: 個別データ情報をシートに書き込む
	// // SND詳細データを更新
	// update_sheet_snd_def(f, snd_data, styleMap)
	// // RCV詳細データを更新
	// update_sheet_rcv_def(f, rcv_data, styleMap)
	// // HST詳細データを更新
	// update_sheet_hst_def(f, hst_data, styleMap)
	// // TGRP詳細データを更新
	// update_sheet_tgrp_def(f, tgrp_data, styleMap)
	// // JOB詳細データを更新
	// update_sheet_job_def(f, job_data, styleMap)
	// // FMT詳細データを更新
	// update_sheet_fmt_def(f, fmt_data, styleMap)
	// // MFMT詳細データを更新
	// update_sheet_mfmt_def(f, mfmt_data, styleMap)
	// // TRG詳細データを更新
	// update_sheet_trg_def(f, trg_data, styleMap)

	// Index
	update_sheet_index(f)

	// 生成されたExcelファイルを指定パスに保存
	// if err := f.SaveAs(makedBook); err != nil {
	if err := f.SaveAs(cfg.MakedBookName); err != nil {
		fmt.Println(err)
	}

}
