package control

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/kngoofy/go-kntool-jp1ajs/pkg/model"
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

func ParseFile(r io.Reader) (*model.UnitBlock, []string) {
	// f, err := os.Open(path)
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()

	scanner := bufio.NewScanner(r)

	root := &model.UnitBlock{UnitName: "root"}
	stack := []*model.UnitBlock{root}

	unitName := ""
	PermissionMode := ""
	Jp1UserName := ""
	Jp1ResourcesGroup := ""
	AbsoluteName := ""
	AbsName := []string{}

	lines := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		trim := strings.TrimSpace(line)
		// ブロック開始
		if strings.Contains(trim, "unit=") {
			// level++
			unitName = strings.Split(strings.Split(trim, "=")[1], ",")[0]
			PermissionMode = strings.Split(strings.Split(trim, "=")[1], ",")[1]
			Jp1UserName = strings.Split(strings.Split(trim, "=")[1], ",")[2]
			Jp1ResourcesGroup = strings.Split(strings.Split(trim, "=")[1], ",")[3]
			Jp1ResourcesGroup = strings.Split(Jp1ResourcesGroup, ":")[0]

			AbsName = append(AbsName, unitName)
			AbsoluteName = "/" + strings.Join(AbsName, "/")

			continue
		}

		// ブロック開始
		if strings.Contains(trim, "{") {
			// level++
			parent := stack[len(stack)-1]
			newBlock := &model.UnitBlock{UnitName: unitName,
				UnitParm: model.UnitParm{
					PermissionMode:    PermissionMode,
					Jp1UserName:       Jp1UserName,
					Jp1ResourcesGroup: Jp1ResourcesGroup},
				ParentUnitName:   parent.UnitName,
				UnitAbsoluteName: AbsoluteName}

			parent.Children = append(parent.Children, newBlock)

			stack = append(stack, newBlock)
			continue
		}

		// ブロック終了
		if strings.Contains(trim, "}") {
			// level--
			stack = stack[:len(stack)-1]
			AbsName = AbsName[:len(stack)-1]
			continue
		}

		// ブロック内の通常行
		if len(stack) > 1 {
			stack[len(stack)-1].Content += trim + "\n"
			stack[len(stack)-1].ContentLine = append(stack[len(stack)-1].ContentLine, trim)
		}
	}

	// fmt.Println("stop")
	return root, lines
}

// 抽象構文木を表示
func PrintAST(b *model.UnitBlock, indent int) {
	prefix := strings.Repeat("  ", indent)
	_ = prefix
	// fmt.Printf("%sBlock: %s\n", prefix, b.Name)
	// fmt.Printf("%sBlock: %s,%s,%s,%s :%s,%s:\n", prefix,
	// 	b.UnitName, b.PermissionMode, b.Jp1UserName, b.Jp1ResourcesGroup, b.UnitAbsoluteName, b.ParentUnitName)
	fmt.Printf("Unit: %s,%s,%s,%s :%s,%s:\n",
		b.UnitName, b.PermissionMode, b.Jp1UserName, b.Jp1ResourcesGroup, b.UnitAbsoluteName, b.ParentUnitName)
	// if strings.TrimSpace(b.Content) != "" {
	// 	fmt.Printf("%s  Content:\n%s%s\n", prefix, prefix, indentText(b.Content, prefix+"    "))
	// }

	for i, line := range b.ContentLine {
		fmt.Println(" :", i, line)
	}

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

func PrintLines(lines []string) {
	for i, r := range lines {
		fmt.Println(i, r)
	}

}

var (
	re_FLCK   = regexp.MustCompile(`FLCK_*`)
	re_ar     = regexp.MustCompile(`ar=\(*`)
	re_arParm = regexp.MustCompile(`ar=\(f=(.+),t=(.+),(.+)\);`)
)

// func init() {
// 	re := regexp.MustCompile(`FLCK_*`)
// }

func FLCKList(b *model.UnitBlock, f *[]*model.UnitBlock) {
	// println(b.UnitName)
	if re_FLCK.MatchString(b.UnitName) {
		// println(b.UnitName)
		*f = append(*f, b)
	}

	for _, child := range b.Children {
		FLCKList(child, f)
	}
}

func ArList(b *model.UnitBlock, f *[]model.ArParm) {

	for _, r := range b.ContentLine {
		if re_ar.MatchString(r) {
			// println(r)
			m := re_arParm.FindStringSubmatch(r)
			arParm := model.ArParm{
				UnitName:         b.UnitName,
				Prior:            m[1],
				Following:        m[2],
				Category:         m[3],
				UnitAbsoluteName: b.UnitAbsoluteName,
			}
			// ar := model.ArParm{}
			*f = append(*f, arParm)
		}
	}

	for _, child := range b.Children {
		ArList(child, f)
	}
}

func JobList(b *model.UnitBlock, j *[]model.JobUnit, n *[]model.NetUnit) {
	isJob := false
	isNet := false
	for _, r := range b.ContentLine {
		if strings.Contains(r, "ty=j") {
			isJob = true
		}
		if strings.Contains(r, "ty=n") {
			isNet = true
		}
	}
	if isJob {
		_jobUnit := model.JobUnit{
			UnitName:         b.UnitName,
			Ty:               "j",
			UnitAbsoluteName: b.UnitAbsoluteName,
		}
		for _, r := range b.ContentLine {
			switch {
			case strings.Contains(r, "cm="):
				_jobUnit.Cm = strings.Split(r, "\"")[1]
			case strings.Contains(r, "ha="):
				_jobUnit.Ha = strings.Split(r, "\"")[1]
			case strings.Contains(r, "te="):
				_jobUnit.Te = strings.Split(r, "\"")[1]
			case strings.Contains(r, "tho="):
				_jobUnit.Tho = strings.Split(r, "=")[1]
				_jobUnit.Tho = _jobUnit.Tho[:len(_jobUnit.Tho)-1]
			case strings.Contains(r, "eu="):
				_jobUnit.Eu = strings.Split(r, "=")[1]
				_jobUnit.Eu = _jobUnit.Eu[:len(_jobUnit.Eu)-1]
			case strings.Contains(r, "un="):
				_jobUnit.Un = strings.Split(r, "\"")[1]
			}

		}
		*j = append(*j, _jobUnit)
	}

	if isNet {
		_netUnit := model.NetUnit{
			UnitName:         b.UnitName,
			Ty:               "n",
			UnitAbsoluteName: b.UnitAbsoluteName,
		}
		for _, r := range b.ContentLine {
			switch {
			case strings.Contains(r, "cm="):
				_netUnit.Cm = strings.Split(r, "\"")[1]
			case strings.Contains(r, "ha="):
				_netUnit.Ha = strings.Split(r, "=")[1]
			case strings.Contains(r, "sz="):
				_netUnit.Sz = strings.Split(r, "=")[1]
				_netUnit.Sz = _netUnit.Sz[:len(_netUnit.Sz)-1]
			case strings.Contains(r, "sd="):
				_netUnit.Sd = strings.Split(r, "=")[1]
				_netUnit.Sd = _netUnit.Sd[:len(_netUnit.Sd)-1]
			case strings.Contains(r, "st="):
				_netUnit.St = strings.Split(r, "=")[1]
				_netUnit.St = _netUnit.St[:len(_netUnit.St)-1]
			case strings.Contains(r, "sh="):
				_netUnit.Sh = strings.Split(r, "=")[1]
				_netUnit.Sh = _netUnit.Sh[:len(_netUnit.Sh)-1]
			case strings.Contains(r, "shd="):
				_netUnit.Shd = strings.Split(r, "=")[1]
				_netUnit.Shd = _netUnit.Shd[:len(_netUnit.Shd)-1]
			case strings.Contains(r, "cy="):
				_netUnit.Cy = strings.Split(r, "=")[1]
				_netUnit.Cy = _netUnit.Cy[:len(_netUnit.Cy)-1]
			case strings.Contains(r, "de="):
				_netUnit.De = strings.Split(r, "=")[1]
				_netUnit.De = _netUnit.De[:len(_netUnit.De)-1]
			}

		}
		*n = append(*n, _netUnit)
	}

	for _, child := range b.Children {
		JobList(child, j, n)
	}
}
