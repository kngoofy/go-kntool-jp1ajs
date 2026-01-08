// =============================================================================
// パッケージ: control
// 概要: JP1/AJSユニット定義ファイルの解析と制御ロジックを提供するパッケージ
// 機能:
//   - ユニット定義ファイルのパース（ParseFile）
//   - 抽象構文木の表示（PrintAST）
//   - ユニット定義行の表示（PrintLines）
//   - FLCKユニットの抽出（FLCKList）
//   - AR（先行・後続関係）パラメータの抽出（ArList）
//   - ジョブ/ネットワークユニットの抽出（JobList）
//
// =============================================================================
package control

import (
	"bufio"   // バッファ付きI/O処理用
	"fmt"     // 書式付き出力用
	"io"      // I/Oプリミティブ用
	"regexp"  // 正規表現処理用
	"strings" // 文字列操作用

	// 内部パッケージ
	"github.com/kngoofy/go-kntool-jp1ajs/pkg/model" // JP1/AJSデータモデル
)

// func ReadFile(r io.Reader) []string {
// 	lines := []string{}

// 	scanner := bufio.NewScanner(r)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		lines = append(lines, line)
// 	}

// 	for i, rec := range lines {
// 		println(i, rec)
// 	}

// 	return lines
// }

// ParseFile はJP1/AJSユニット定義ファイルを解析し、抽象構文木を構築する
// 引数:
//   - r: io.Reader（ユニット定義ファイルの入力ストリーム）
//
// 戻り値:
//   - *model.UnitBlock: 抽象構文木のルートノード
//   - []string: ユニット定義ファイルの各行を格納したスライス
//
// 処理概要:
//  1. ファイルを行単位で読み込む
//  2. "unit=" 行でユニット定義の開始を検知し、ユニット名やパラメータを抽出
//  3. "{" でブロック開始、"}" でブロック終了として階層構造を構築
//  4. ブロック内の通常行はContentLineに格納
func ParseFile(r io.Reader) (*model.UnitBlock, []string) {
	// ※ 旧コード: ファイルパスから直接オープンする処理（現在は呼び出し元で実施）
	// f, err := os.Open(path)
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()

	// ファイルを行単位で読み込むためのスキャナーを作成
	scanner := bufio.NewScanner(r)

	// 抽象構文木の初期化
	// root: ルートノード（全ユニットの親となる仮想ノード）
	// stack: 現在の階層を追跡するためのスタック
	root := &model.UnitBlock{UnitName: "root"}
	stack := []*model.UnitBlock{root}

	// ユニットパラメータ用の一時変数
	unitName := ""          // ユニット名
	PermissionMode := ""    // パーミッションモード
	Jp1UserName := ""       // JP1ユーザー名
	Jp1ResourcesGroup := "" // JP1リソースグループ
	AbsoluteName := ""      // ユニットの絶対パス名
	AbsName := []string{}   // 絶対パス構築用のスライス

	// 読み込んだ行を保持するスライス
	lines := []string{}

	// =========================================================================
	// メインループ: ファイルを1行ずつ読み込んで解析
	// =========================================================================
	for scanner.Scan() {
		// 現在行を取得し、保存用スライスに追加
		line := scanner.Text()
		lines = append(lines, line)

		// 前後の空白を除去した文字列で判定
		trim := strings.TrimSpace(line)

		// -----------------------------------------------------------------
		// ユニット定義行の検出: "unit=" を含む行
		// 形式: unit=ユニット名,パーミッション,JP1ユーザー,リソースグループ:...
		// -----------------------------------------------------------------
		if strings.Contains(trim, "unit=") {
			// ユニット名を抽出（"="の後、最初の","まで）
			unitName = strings.Split(strings.Split(trim, "=")[1], ",")[0]
			// パーミッションモードを抽出
			PermissionMode = strings.Split(strings.Split(trim, "=")[1], ",")[1]
			// JP1ユーザー名を抽出
			Jp1UserName = strings.Split(strings.Split(trim, "=")[1], ",")[2]
			// JP1リソースグループを抽出（":"より前の部分）
			Jp1ResourcesGroup = strings.Split(strings.Split(trim, "=")[1], ",")[3]
			Jp1ResourcesGroup = strings.Split(Jp1ResourcesGroup, ":")[0]

			// 絶対パス名を構築（例: /親ユニット/子ユニット）
			AbsName = append(AbsName, unitName)
			AbsoluteName = "/" + strings.Join(AbsName, "/")

			continue
		}

		// -----------------------------------------------------------------
		// ブロック開始: "{" を含む行
		// 新しいUnitBlockを作成し、親ノードの子として追加
		// -----------------------------------------------------------------
		if strings.Contains(trim, "{") {
			// 現在のスタックトップを親ノードとして取得
			parent := stack[len(stack)-1]

			// 新しいユニットブロックを作成
			newBlock := &model.UnitBlock{UnitName: unitName,
				UnitParm: model.UnitParm{
					PermissionMode:    PermissionMode,
					Jp1UserName:       Jp1UserName,
					Jp1ResourcesGroup: Jp1ResourcesGroup},
				ParentUnitName:   parent.UnitName,
				UnitAbsoluteName: AbsoluteName}

			// 親ノードの子リストに追加
			parent.Children = append(parent.Children, newBlock)

			// スタックにプッシュ（次の行からはこのブロック内として処理）
			stack = append(stack, newBlock)
			continue
		}

		// -----------------------------------------------------------------
		// ブロック終了: "}" を含む行
		// スタックから現在のブロックをポップし、階層を一つ上に戻る
		// -----------------------------------------------------------------
		if strings.Contains(trim, "}") {
			// スタックからポップ（現在のブロックを終了）
			stack = stack[:len(stack)-1]
			// 絶対パス名も合わせて調整
			AbsName = AbsName[:len(stack)-1]
			continue
		}

		// -----------------------------------------------------------------
		// ブロック内の通常行（パラメータ定義など）
		// 現在のユニットブロックのContentとContentLineに追加
		// -----------------------------------------------------------------
		if len(stack) > 1 {
			// Content: 改行区切りの連結文字列として保持
			stack[len(stack)-1].Content += trim + "\n"
			// ContentLine: 各行を個別にスライスとして保持
			stack[len(stack)-1].ContentLine = append(stack[len(stack)-1].ContentLine, trim)
		}
	}

	// 構築した抽象構文木のルートノードと読み込んだ行を返却
	return root, lines
}

// PrintAST は抽象構文木を再帰的に標準出力に表示する
// 引数:
//   - b: 表示対象のUnitBlockノード
//   - indent: 現在のインデントレベル（階層の深さ）
//
// 処理概要:
//  1. ユニット情報（名前、パーミッション、ユーザー、リソースグループ等）を表示
//  2. ユニット内の各コンテンツ行を表示
//  3. 子ノードに対して再帰的に同じ処理を実行
func PrintAST(b *model.UnitBlock, indent int) {
	// インデント用のプレフィックス文字列を生成（現在は未使用）
	prefix := strings.Repeat("  ", indent)
	_ = prefix // 将来の拡張用に保持

	// ユニット情報を1行で出力
	// 形式: Unit: ユニット名,パーミッション,JP1ユーザー,リソースグループ :絶対パス,親ユニット名:
	fmt.Printf("Unit: %s,%s,%s,%s :%s,%s:\n",
		b.UnitName, b.PermissionMode, b.Jp1UserName, b.Jp1ResourcesGroup, b.UnitAbsoluteName, b.ParentUnitName)

	// ユニット内のコンテンツ行を番号付きで出力
	for i, line := range b.ContentLine {
		fmt.Println(" :", i, line)
	}

	// 子ノードに対して再帰呼び出し（階層を1つ深くする）
	for _, child := range b.Children {
		PrintAST(child, indent+1)
	}
}

//	func indentText(text, indent string) string {
//		lines := strings.Split(text, "\n")
//		for i := range lines {
//			lines[i] = indent + lines[i]
//		}
//		return strings.Join(lines, "\n")
//	}

// PrintLines はユニット定義ファイルの各行を行番号付きで表示する
// 引数:
//   - lines: 表示対象の行スライス
//
// 用途: デバッグや確認用に、読み込んだファイル内容を確認する
func PrintLines(lines []string) {
	// 各行を行番号と共に出力
	for i, r := range lines {
		fmt.Println(i, r)
	}
}

// =============================================================================
// 正規表現パターン定義（パッケージレベルで事前コンパイル）
// =============================================================================
var (
	// re_FLCK: FLCKユニット名にマッチする正規表現（FLCK_で始まるユニット）
	// re_FLCK = regexp.MustCompile(`FLCK_*`)
	// re_ar: AR（先行・後続）パラメータの存在チェック用
	re_ar = regexp.MustCompile(`ar=\(*`)
	// re_arParm: ARパラメータの詳細抽出用
	// 形式: ar=(f=先行ユニット,t=後続ユニット,関係種別);
	// キャプチャグループ: 1=先行, 2=後続, 3=種別
	re_arParm = regexp.MustCompile(`ar=\(f=(.+),t=(.+),(.+)\);`)
)

// func init() {
// 	re := regexp.MustCompile(`FLCK_*`)
// }

// FLCKList は抽象構文木からFLCK（ファイルロック）ユニットを抽出する
// 引数:
//   - b: 検索対象のUnitBlockノード
//   - f: 抽出結果を格納するスライスへのポインタ
//
// 処理概要:
//  1. ユニット名がFLCK_で始まる場合、結果スライスに追加
//  2. 子ノードに対して再帰的に同じ処理を実行
func FLCKList(b *model.UnitBlock, f *[]model.FlckUnit) {

	isFlck := false

	// ユニット種別を判定
	for _, r := range b.ContentLine {
		if strings.Contains(r, "ty=flwj") {
			isFlck = true
		}
	}

	if isFlck {
		_flckUnit := model.FlckUnit{
			UnitName:         b.UnitName,
			Ty:               "jflwj",
			UnitAbsoluteName: b.UnitAbsoluteName,
		}

		// _ = _flckUnit

		for _, r := range b.ContentLine {
			switch {
			case strings.Contains(r, "cm="):
				_flckUnit.Cm = strings.Split(r, "\"")[1]
			case strings.Contains(r, "flwf="):
				_flckUnit.Flwf = strings.Split(r, "\"")[1]
			case strings.Contains(r, "flwc="):
				_flckUnit.Flwc = strings.Split(r, "=")[1]
				_flckUnit.Flwc = _flckUnit.Flwc[:len(_flckUnit.Flwc)-1]
			case strings.Contains(r, "flco="):
				_flckUnit.Flco = strings.Split(r, "=")[1]
				_flckUnit.Flco = _flckUnit.Flco[:len(_flckUnit.Flco)-1]
			case strings.Contains(r, "flwi="):
				_flckUnit.Flwi = strings.Split(r, "=")[1]
				_flckUnit.Flwi = _flckUnit.Flwi[:len(_flckUnit.Flwi)-1]
			case strings.Contains(r, "eu="):
				_flckUnit.Eu = strings.Split(r, "=")[1]
				_flckUnit.Eu = _flckUnit.Eu[:len(_flckUnit.Eu)-1]
			}
		}

		*f = append(*f, _flckUnit)
	}

	// if re_FLCK.MatchString(b.UnitName) {
	// 	// マッチした場合、結果スライスに追加
	// }

	// 子ノードに対して再帰呼び出し
	for _, child := range b.Children {
		FLCKList(child, f)
	}
}

// ArList は抽象構文木からAR（先行・後続関係）パラメータを抽出する
// 引数:
//   - b: 検索対象のUnitBlockノード
//   - f: 抽出結果を格納するスライスへのポインタ
//
// 処理概要:
//  1. ユニットのContentLineから "ar=" を含む行を検索
//  2. 正規表現で先行ユニット、後続ユニット、関係種別を抽出
//  3. ArParm構造体を作成し、結果スライスに追加
//  4. 子ノードに対して再帰的に同じ処理を実行
func ArList(b *model.UnitBlock, f *[]model.ArParm) {
	// ユニット内の各コンテンツ行を走査
	for _, r := range b.ContentLine {
		// AR定義が存在するかチェック
		if re_ar.MatchString(r) {
			// 正規表現でARパラメータの詳細を抽出
			// m[0]=マッチ全体, m[1]=先行, m[2]=後続, m[3]=種別
			m := re_arParm.FindStringSubmatch(r)

			// ArParm構造体を作成
			arParm := model.ArParm{
				UnitName:         b.UnitName,         // 所属ユニット名
				Prior:            m[1],               // 先行ユニット
				Following:        m[2],               // 後続ユニット
				Category:         m[3],               // 関係種別
				UnitAbsoluteName: b.UnitAbsoluteName, // ユニットの絶対パス
			}

			// 結果スライスに追加
			*f = append(*f, arParm)
		}
	}

	// 子ノードに対して再帰呼び出し
	for _, child := range b.Children {
		ArList(child, f)
	}
}

// JobList は抽象構文木からジョブユニットとネットワークユニットを抽出する
// 引数:
//   - b: 検索対象のUnitBlockノード
//   - j: ジョブユニットを格納するスライスへのポインタ
//   - n: ネットワークユニットを格納するスライスへのポインタ
//
// 処理概要:
//  1. ユニットの種別（ty=j: ジョブ, ty=n: ネットワーク）を判定
//  2. ジョブの場合: cm, ha, te, tho, eu, un パラメータを抽出
//  3. ネットワークの場合: cm, ha, sz, sd, st, sh, shd, cy, de パラメータを抽出
//  4. 子ノードに対して再帰的に同じ処理を実行
func JobList(b *model.UnitBlock, j *[]model.JobUnit, n *[]model.NetUnit) {
	// ユニット種別判定フラグ
	isJob := false // ty=j (ジョブユニット)
	isNet := false // ty=n (ネットワークユニット)

	// ユニット種別を判定
	for _, r := range b.ContentLine {
		if strings.Contains(r, "ty=j") {
			isJob = true
		}
		if strings.Contains(r, "ty=n") {
			isNet = true
		}
	}

	// =========================================================================
	// ジョブユニット（ty=j）の処理
	// =========================================================================
	if isJob {
		// JobUnit構造体を初期化
		_jobUnit := model.JobUnit{
			UnitName:         b.UnitName,
			Ty:               "j", // ユニット種別: ジョブ
			UnitAbsoluteName: b.UnitAbsoluteName,
		}

		// 各パラメータを抽出
		for _, r := range b.ContentLine {
			switch {
			case strings.Contains(r, "cm="):
				// cm: コメント（ダブルクォート内の値）
				_jobUnit.Cm = strings.Split(r, "\"")[1]
			case strings.Contains(r, "ha="):
				// ha: ホストエージェント（ダブルクォート内の値）
				_jobUnit.Ha = strings.Split(r, "\"")[1]
			case strings.Contains(r, "te="):
				// te: 実行ファイル名（ダブルクォート内の値）
				_jobUnit.Te = strings.Split(r, "\"")[1]
			case strings.Contains(r, "tho="):
				// tho: タイムアウト時間（末尾のセミコロンを除去）
				_jobUnit.Tho = strings.Split(r, "=")[1]
				_jobUnit.Tho = _jobUnit.Tho[:len(_jobUnit.Tho)-1]
			case strings.Contains(r, "eu="):
				// eu: 実行ユーザー種別（末尾のセミコロンを除去）
				_jobUnit.Eu = strings.Split(r, "=")[1]
				_jobUnit.Eu = _jobUnit.Eu[:len(_jobUnit.Eu)-1]
			case strings.Contains(r, "un="):
				// un: 実行ユーザー名（ダブルクォート内の値）
				_jobUnit.Un = strings.Split(r, "\"")[1]
			}
		}

		// 結果スライスに追加
		*j = append(*j, _jobUnit)
	}

	// =========================================================================
	// ネットワークユニット（ty=n）の処理
	// =========================================================================
	if isNet {
		// NetUnit構造体を初期化
		_netUnit := model.NetUnit{
			UnitName:         b.UnitName,
			Ty:               "n", // ユニット種別: ネットワーク
			UnitAbsoluteName: b.UnitAbsoluteName,
		}

		// 各パラメータを抽出
		for _, r := range b.ContentLine {
			switch {
			case strings.Contains(r, "cm="):
				// cm: コメント（ダブルクォート内の値）
				_netUnit.Cm = strings.Split(r, "\"")[1]
			case strings.Contains(r, "ha="):
				// ha: ホストエージェント
				_netUnit.Ha = strings.Split(r, "=")[1]
			case strings.Contains(r, "sz="):
				// sz: サイズ（末尾のセミコロンを除去）
				_netUnit.Sz = strings.Split(r, "=")[1]
				_netUnit.Sz = _netUnit.Sz[:len(_netUnit.Sz)-1]
			case strings.Contains(r, "sd="):
				// sd: 開始日（末尾のセミコロンを除去）
				_netUnit.Sd = strings.Split(r, "=")[1]
				_netUnit.Sd = _netUnit.Sd[:len(_netUnit.Sd)-1]
			case strings.Contains(r, "st="):
				// st: 開始時刻（末尾のセミコロンを除去）
				_netUnit.St = strings.Split(r, "=")[1]
				_netUnit.St = _netUnit.St[:len(_netUnit.St)-1]
			case strings.Contains(r, "sh="):
				// sh: スケジュールホスト（末尾のセミコロンを除去）
				_netUnit.Sh = strings.Split(r, "=")[1]
				_netUnit.Sh = _netUnit.Sh[:len(_netUnit.Sh)-1]
			case strings.Contains(r, "shd="):
				// shd: スケジュール日（末尾のセミコロンを除去）
				_netUnit.Shd = strings.Split(r, "=")[1]
				_netUnit.Shd = _netUnit.Shd[:len(_netUnit.Shd)-1]
			case strings.Contains(r, "cy="):
				// cy: 実行サイクル（末尾のセミコロンを除去）
				_netUnit.Cy = strings.Split(r, "=")[1]
				_netUnit.Cy = _netUnit.Cy[:len(_netUnit.Cy)-1]
			case strings.Contains(r, "de="):
				// de: 遅延監視（末尾のセミコロンを除去）
				_netUnit.De = strings.Split(r, "=")[1]
				_netUnit.De = _netUnit.De[:len(_netUnit.De)-1]
			}
		}

		// 結果スライスに追加
		*n = append(*n, _netUnit)
	}

	// 子ノードに対して再帰呼び出し
	for _, child := range b.Children {
		JobList(child, j, n)
	}
}
