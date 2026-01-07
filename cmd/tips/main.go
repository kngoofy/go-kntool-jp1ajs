package main

import (
	"fmt"
	"regexp"

	"github.com/kngoofy/go-kntool-jp1ajs/pkg/model"
)

func main() {
	str := "ar=(f=jobnet01,t=jobnet02,seq);"

	re := regexp.MustCompile(`ar=\(f=(.+),t=(.+),(.+)\);`)
	m := re.FindStringSubmatch(str)
	prior := m[1]
	following := m[2]
	category := m[3]

	fmt.Println(prior, following, category)

	ArParm := model.ArParm{
		Prior:     m[1],
		Following: m[2],
		Category:  m[3],
	}

	fmt.Println(ArParm)
	// _ = ArParm

	// re2 := regexp.MustCompile(`f=(.+),`)
	// m2 := re2.FindStringSubmatch(str)

	// re1 := regexp.MustCompile(`\(.+\)`)
	// m1 := re1.FindStringSubmatch(str)
	// reBefore := regexp.MustCompile(`\(f=(.+),t=`)
	// before := reBefore.FindStringSubmatch(m1[0])[1]
	// reBefore = regexp.MustCompile(`ar=\(f=(.+),t=(.+),(.+)\);`)
	// beforeA := reBefore.FindStringSubmatch(str)

	// reAfter := regexp.MustCompile(`t=(.+),`)
	// after := reAfter.FindStringSubmatch(str)[1]

	// reSeq := regexp.MustCompile(`,(.+)\)`)
	// seq := reSeq.FindStringSubmatch(str)

	// fmt.Println(prior, following, category)

	// re5 := regexp.MustCompile(`,+`)
	// m5 := re5.Split(str, -1)
	// _ = m
	// _ = m5
	// _ = m1
	// _ = m2
	// _ = after
	// _ = before
	// _ = seq
	// _ = beforeA

	if len(m) > 1 {
		fmt.Println(m[1]) // echo test
	}

}
