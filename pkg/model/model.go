// package main は、JP1/AJS 定義などをパースするためのデータモデルを定義します。
package model

type UnitParm struct {
	PermissionMode    string
	Jp1UserName       string
	Jp1ResourcesGroup string
}
type UnitBlock struct {
	UnitName string
	UnitParm
	UnitAbsoluteName string
	Content          string
	ContentLine      []string
	Children         []*UnitBlock
	ParentUnitName   string
	// PermissionMode    string
	// Jp1UserName       string
	// Jp1ResourcesGroup string
}

type ArParm struct {
	UnitName         string
	Prior            string
	Following        string
	Category         string
	UnitAbsoluteName string
}

type NetUnit struct {
	UnitName         string
	Sz               string
	Ty               string
	Cm               string
	Ha               string
	Sd               string
	St               string
	Cy               string
	Sh               string
	Shd              string
	De               string
	UnitAbsoluteName string
}

type JobUnit struct {
	UnitName         string
	Ty               string
	Cm               string
	Ha               string
	Te               string
	Tho              string
	Eu               string
	Un               string
	UnitAbsoluteName string
}

// Unit は、全てのユニット（ジョブ、ネット、グループ等）のベースとなる共通構造体です。
type Unit struct {
	UnitName       string   // ユニット名
	UnitAbsName    string   // フルパス名（絶対名）
	UnitType       string   // ユニット種別（ty=...）
	Comment        string   // コメント（cm=...）
	ParentUnitName string   // 親ユニット名
	ChildUnitName  []string // 子ユニット名のリスト
	ChildPointer   *any     // 子要素へのポインタ（拡張用）
}

// // グループUnit
// type UnitGroup struct {
// 	Unit
// }

// ネットUnit
type UnitNet struct {
	Unit
	El
	Ar
	Sz
}

// ジョブUnit
type UnitJob struct {
	Unit
}

// sz
type Sz struct {
	Size string
}

// el
type El struct {
	UnitName string
	Location string
}

// el
type Ar struct {
	f string
	t string
	n string
}
