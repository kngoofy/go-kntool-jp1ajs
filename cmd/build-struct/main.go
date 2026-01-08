// =============================================================================
// パッケージ: main
// 概要: JP1/AJS（Hitachi Job Management Partner 1/Automatic Job Management System）の
//
//	ユニット定義ファイルを解析し、抽象構文木（AST）を構築するツール
//
// 機能:
//   - ユニット定義ファイルの読み込みとパース
//   - 抽象構文木の構築と表示
//   - FLCK（ファイルロック）リストの抽出
//   - AR（関連パラメータ）リストの抽出
//   - ジョブユニットおよびネットワークユニットリストの抽出
//
// =============================================================================
package main

import (
	"fmt" // 標準出力用パッケージ
	"os"  // ファイル操作用パッケージ

	// 内部パッケージ
	"github.com/kngoofy/go-kntool-jp1ajs/pkg/control" // JP1/AJS制御ロジック
	"github.com/kngoofy/go-kntool-jp1ajs/pkg/model"   // JP1/AJSデータモデル
)

// type UnitParm struct {
// 	PermissionMode    string
// 	Jp1UserName       string
// 	Jp1ResourcesGroup string
// }
// type UnitBlock struct {
// 	UnitName string
// 	UnitParm
// 	UnitAbsoluteName string
// 	Content          string
// 	ContentLine      []string
// 	Children         []*UnitBlock
// 	ParentUnitName   string
// 	// PermissionMode    string
// 	// Jp1UserName       string
// 	// Jp1ResourcesGroup string
// }

// main関数: プログラムのエントリーポイント
// JP1/AJSユニット定義ファイルを解析し、各種情報を抽出・表示する
func main() {

	// =========================================================================
	// [1].JP1/AJSのユニット定義ファイルをOpen
	// =========================================================================
	// 解析対象のユニット定義ファイルパスを指定
	// file := "../../pkg/testdata/input2.def" // 代替ファイル（コメントアウト）
	file := "../../pkg/testdata/jp1ajs.txt" // 現在の解析対象ファイル

	// ファイルを読み取り専用モードでオープン
	f, err := os.Open(file)
	if err != nil {
		// ファイルオープン失敗時はパニックで終了
		panic(err)
	}
	// 関数終了時にファイルを確実にクローズ
	defer f.Close()

	// =========================================================================
	// [2].ユニット定義をパースし抽象構文木を組み立てる
	// =========================================================================
	// ParseFile: ファイルを行単位で読み込み、ユニット定義の階層構造を解析
	// 戻り値:
	//   - root: 抽象構文木のルートノード（UnitBlock型）
	//   - lines: ユニット定義ファイルの各行を格納したスライス（[]string型）
	root, lines := control.ParseFile(f)
	_ = root  // 未使用警告を抑制（後続処理で使用）
	_ = lines // 未使用警告を抑制（後続処理で使用）

	// =========================================================================
	// [3].パース結果を表示（抽象構文木）
	// =========================================================================
	// 構築された抽象構文木をツリー形式で標準出力に表示
	// インデントレベル0から開始し、再帰的に子ノードを表示
	fmt.Println("\n- 抽象構文木 -")
	control.PrintAST(root, 0)

	// =========================================================================
	// [4].読込み組み立てたユニット定義ファイルを表示
	// =========================================================================
	// パース時に取得した各行を順番に表示（デバッグ・確認用）
	fmt.Println("\n- JP1/AJSのユニット定義ファイル -")
	control.PrintLines(lines)

	// =========================================================================
	// [5].FLCK（ファイルロック）リストの抽出
	// =========================================================================
	// FLCKList: 抽象構文木からファイルロック関連のユニットを抽出
	// flckStack: 抽出されたFLCKユニットを格納するスライス（ポインタ渡し）
	flckStack := []model.FlckUnit{}
	control.FLCKList(root, &flckStack)

	// =========================================================================
	// [6].AR（関連パラメータ）リストの抽出
	// =========================================================================
	// ArList: 抽象構文木からAR（先行・後続関係）パラメータを抽出
	// arStack: 抽出されたARパラメータを格納するスライス（ポインタ渡し）
	arStack := []model.ArParm{}
	control.ArList(root, &arStack)

	// =========================================================================
	// [7].ジョブユニットおよびネットワークユニットリストの抽出
	// =========================================================================
	// JobList: 抽象構文木からジョブユニットとネットワークユニットを分類・抽出
	// jobStack: 抽出されたジョブユニットを格納するスライス
	// netStack: 抽出されたネットワークユニットを格納するスライス
	jobStack := []model.JobUnit{}
	netStack := []model.NetUnit{}
	control.JobList(root, &jobStack, &netStack)

}

// func parseFile(path string) *UnitBlock {
// 	f, err := os.Open(path)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()

// 	scanner := bufio.NewScanner(f)

// 	root := &UnitBlock{UnitName: "root"}
// 	stack := []*UnitBlock{root}

// 	unitName := ""
// 	PermissionMode := ""
// 	Jp1UserName := ""
// 	Jp1ResourcesGroup := ""
// 	AbsoluteName := ""
// 	AbsName := []string{}

// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		trim := strings.TrimSpace(line)
// 		// ブロック開始
// 		if strings.Contains(trim, "unit=") {
// 			// level++
// 			unitName = strings.Split(strings.Split(trim, "=")[1], ",")[0]
// 			PermissionMode = strings.Split(strings.Split(trim, "=")[1], ",")[1]
// 			Jp1UserName = strings.Split(strings.Split(trim, "=")[1], ",")[2]
// 			Jp1ResourcesGroup = strings.Split(strings.Split(trim, "=")[1], ",")[3]
// 			Jp1ResourcesGroup = strings.Split(Jp1ResourcesGroup, ":")[0]

// 			AbsName = append(AbsName, unitName)
// 			AbsoluteName = strings.Join(AbsName, "/")

// 			continue
// 		}

// 		// ブロック開始
// 		if strings.Contains(trim, "{") {
// 			// level++
// 			parent := stack[len(stack)-1]
// 			newBlock := &UnitBlock{UnitName: unitName,
// 				UnitParm: UnitParm{
// 					PermissionMode:    PermissionMode,
// 					Jp1UserName:       Jp1UserName,
// 					Jp1ResourcesGroup: Jp1ResourcesGroup},
// 				ParentUnitName:   parent.UnitName,
// 				UnitAbsoluteName: AbsoluteName}

// 			parent.Children = append(parent.Children, newBlock)

// 			stack = append(stack, newBlock)
// 			continue
// 		}

// 		// ブロック終了
// 		if strings.Contains(trim, "}") {
// 			// level--
// 			stack = stack[:len(stack)-1]
// 			AbsName = AbsName[:len(stack)-1]
// 			continue
// 		}

// 		// ブロック内の通常行
// 		if len(stack) > 1 {
// 			stack[len(stack)-1].Content += trim + "\n"
// 			stack[len(stack)-1].ContentLine = append(stack[len(stack)-1].ContentLine, trim)
// 		}
// 	}

// 	// fmt.Println("stop")
// 	return root
// }

// func printAST(b *UnitBlock, indent int) {
// 	prefix := strings.Repeat("  ", indent)
// 	// fmt.Printf("%sBlock: %s\n", prefix, b.Name)
// 	fmt.Printf("%sBlock: %s,%s,%s,%s :%s,%s:\n", prefix,
// 		b.UnitName, b.PermissionMode, b.Jp1UserName, b.Jp1ResourcesGroup, b.UnitAbsoluteName, b.ParentUnitName)
// 	if strings.TrimSpace(b.Content) != "" {
// 		fmt.Printf("%s  Content:\n%s%s\n", prefix, prefix, indentText(b.Content, prefix+"    "))
// 	}

// 	for i, line := range b.ContentLine {
// 		fmt.Println(" :", i, line)
// 	}

// 	for _, child := range b.Children {
// 		printAST(child, indent+1)
// 	}
// }

// func indentText(text, indent string) string {
// 	lines := strings.Split(text, "\n")
// 	for i := range lines {
// 		lines[i] = indent + lines[i]
// 	}
// 	return strings.Join(lines, "\n")
// }
