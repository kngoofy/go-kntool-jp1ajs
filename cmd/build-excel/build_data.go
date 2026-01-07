package main

import (
	"os"

	"github.com/kngoofy/go-kntool-jp1ajs/pkg/control"
	"github.com/kngoofy/go-kntool-jp1ajs/pkg/model"
)

func build_jp1ajs_model(file string) (*model.UnitBlock, []string) {

	// [1].JP1/AJSのユニット定義ファイルをOpen
	// file := "../../pkg/testdata/input2.def"
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// [2].ユニット定義をパースし抽象構文木を組み立てる,ユニット定義の[]string組み立て
	root, lines := control.ParseFile(f)
	_ = root
	_ = lines

	// // [3].パース結果を表示
	// fmt.Println("\n- 抽象構文木 -")
	// control.PrintAST(root, 0)

	// // [4].読込み組み立てたユニット定義ファイルを表示
	// fmt.Println("\n- JP1/AJSのユニット定義ファイル -")
	// control.PrintLines(lines)

	return root, lines
}
