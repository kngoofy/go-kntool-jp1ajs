package main

import (
	"fmt"
	"os"

	"github.com/kngoofy/go-kntool-jp1ajs/pkg/control"
	"github.com/kngoofy/go-kntool-jp1ajs/pkg/model"
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

func main() {

	// [1].JP1/AJSのユニット定義ファイルをOpen
	// file := "../../pkg/testdata/input2.def"
	file := "../../pkg/testdata/jp1ajs.txt"
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// [2].ユニット定義をパースし抽象構文木を組み立てる,ユニット定義の[]string組み立て
	root, lines := control.ParseFile(f)
	_ = root
	_ = lines

	// [3].パース結果を表示
	fmt.Println("\n- 抽象構文木 -")
	control.PrintAST(root, 0)

	// [4].読込み組み立てたユニット定義ファイルを表示
	fmt.Println("\n- JP1/AJSのユニット定義ファイル -")
	control.PrintLines(lines)

	// [5]. s
	flckStack := []*model.UnitBlock{}
	control.FLCKList(root, &flckStack)

	// [6]. s
	arStack := []model.ArParm{}
	control.ArList(root, &arStack)

	// [7]. s
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
