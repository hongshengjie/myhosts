package ui

import (
	"myhosts/model"

	"gioui.org/layout"
	"gioui.org/widget"
)

type Tabs struct {
	list     layout.List
	tabs     []*Tab
	selected int
}

type Tab struct {
	btn     *widget.Clickable
	editor  *widget.Editor
	aswitch *widget.Bool
	host    *model.Host
	saveBtn *widget.Clickable
}
