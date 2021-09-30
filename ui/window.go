package ui

import (
	"image"
	"myhosts/hosts"
	"myhosts/model"
	"myhosts/ui/assets"
	"runtime"

	"gioui.org/app"
	"gioui.org/font/opentype"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gen2brain/beeep"
	"github.com/gen2brain/dlgs"
)

type MainWindow struct {
	theme        *material.Theme
	invalidate   chan struct{}
	ops          *op.Ops
	saveButton   *widget.Clickable
	hostsManager *hosts.HostsFileManager

	//
	split *Split
	//
	tabs *Tabs
	plus *widget.Clickable
	mins *widget.Clickable
}

func CreateWindow() (*MainWindow, *app.Window) {

	ttc, _ := opentype.Parse(assets.Fonts)

	appw := app.NewWindow(app.Title("myhosts"))
	w := &MainWindow{
		invalidate:   make(chan struct{}, 1),
		theme:        material.NewTheme([]text.FontFace{{Font: text.Font{Typeface: "msyh"}, Face: ttc}}),
		ops:          &op.Ops{},
		saveButton:   &widget.Clickable{},
		hostsManager: hosts.NewHostFileManager(),

		split: &Split{Ratio: -0.5},
		plus:  &widget.Clickable{},
		mins:  &widget.Clickable{},
	}
	backup := w.hostsManager.CurrentHostFile()
	if model.FirstOpen() {
		w.hostsManager.Create(&model.Host{
			ID:      0,
			Title:   "backup",
			Content: backup,
			Enable:  true,
		})
		w.hostsManager.ReLoad()
	}

	w.flash()

	return w, appw
}

func (m *MainWindow) flash() {

	tabs := &Tabs{
		list: layout.List{
			Axis: layout.Vertical,
		},
	}
	m.tabs = tabs
	m.tabs.tabs = append(m.tabs.tabs, &Tab{
		btn:  &widget.Clickable{},
		host: &model.Host{Title: "system hosts"},
	})
	for _, v := range m.hostsManager.All() {
		t := &Tab{
			btn:     &widget.Clickable{},
			host:    v,
			editor:  &widget.Editor{},
			aswitch: &widget.Bool{Value: v.Enable},
			saveBtn: &widget.Clickable{},
		}
		t.editor.SetText(v.Content)
		m.tabs.tabs = append(m.tabs.tabs, t)

	}

}

func (m *MainWindow) Loop(w *app.Window, shutdown chan int) error {

	for {
		select {
		case msg := <-m.hostsManager.Err():
			beeep.Notify("错误", msg, "")
			m.hostsManager.SetPwd("")
		case <-m.invalidate:
			m.hostsManager.ReLoad()
			m.flash()
			w.Invalidate()
		case e := <-w.Events():
			switch evt := e.(type) {
			case system.DestroyEvent:
				close(shutdown)
			case system.FrameEvent:
				gtx := layout.NewContext(m.ops, evt)
				m.Action(gtx)
				m.Layout(gtx, m.theme)
				evt.Frame(gtx.Ops)
			default:
			}
		}

	}
}

func (m *MainWindow) Layout(gtx C, th *material.Theme) D {
	return layout.Stack{
		Alignment: layout.NW,
	}.Layout(gtx,
		layout.Expanded(
			func(gtx C) D {
				return m.drawPage(gtx, th)
			},
		),
	)
}

func (m *MainWindow) left(gtx C, th *material.Theme) func(gtx C) D {
	return func(gtx C) D {

		d := layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		}.Layout(gtx,

			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return m.tabs.list.Layout(gtx, len(m.tabs.tabs), func(gtx layout.Context, index int) layout.Dimensions {
					if m.tabs.tabs[index].btn.Clicked() {
						m.tabs.selected = index
					}
					var tabHeight int
					return layout.Stack{Alignment: layout.E}.Layout(gtx,
						layout.Stacked(func(gtx C) D {
							left := func(gtx C) D {
								return layout.UniformInset(unit.Sp(12)).Layout(gtx, func(gtx C) D {
									return material.H6(th, m.tabs.tabs[index].host.Title).Layout(gtx)
								})

							}
							right := func(gtx C) D {
								if index == 0 {
									return D{}
								}
								sw := m.tabs.tabs[index].aswitch
								return layout.UniformInset(unit.Sp(12)).Layout(gtx, func(gtx C) D { return material.Switch(th, sw).Layout(gtx) })
							}

							dims := material.Clickable(gtx, m.tabs.tabs[index].btn, func(gtx C) D {
								return endToEndRow(gtx, left, right)
							})
							tabHeight = dims.Size.Y

							return dims
						}),

						layout.Stacked(func(gtx C) D {
							if m.tabs.selected != index {
								return D{}
							}
							tabWidth := gtx.Px(unit.Dp(4))
							tabRect := image.Rect(0, 0, tabWidth, tabHeight)
							paint.FillShape(gtx.Ops, th.Palette.ContrastBg, clip.Rect(tabRect).Op())
							return D{
								Size: image.Point{X: tabWidth, Y: tabHeight},
							}

						}),
					)

				})
			}),
		)
		return d
	}

}

func (m *MainWindow) right(gtx C, th *material.Theme) func(gtx C) D {
	return func(gtx C) D {
		return layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			layout.Flexed(1,
				func(gtx C) D {
					index := m.tabs.selected
					if index == 0 {
						return material.Label(th, unit.Dp(20), m.hostsManager.CurrentHostFile()).Layout(gtx)
					}
					list := m.tabs.tabs
					if len(list) > index {
						e := list[index].editor
						return material.Editor(th, e, "#example\n127.0.0.1 localhost").Layout(gtx)
					}
					return D{}

				},
			),

			layout.Rigid(func(gtx C) D {
				if m.tabs.selected == 0 {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(material.IconButton(th, m.mins, contentRemove).Layout),
						layout.Rigid(material.IconButton(th, m.plus, contentAdd).Layout),
					)
				}
				save := m.tabs.tabs[m.tabs.selected].saveBtn

				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(material.IconButton(th, m.mins, contentRemove).Layout),
					layout.Rigid(material.IconButton(th, m.plus, contentAdd).Layout),
					layout.Flexed(1, func(gtx C) D {
						return layout.E.Layout(gtx, material.IconButton(th, save, contentSave).Layout)
					}),
				)

			}),
		)

	}

}

func (m *MainWindow) drawPage(gtx C, th *material.Theme) D {
	return m.split.Layout(gtx, th, m.left(gtx, th), m.right(gtx, th))
}

func (m *MainWindow) Action(gtx C) {
	//var changeText string
	var switchChanged, editorSaved bool
	var updateHost []*model.Host
	for _, v := range m.tabs.tabs {

		if v.aswitch != nil && v.aswitch.Changed() {
			switchChanged = true
			v.host.Enable = v.aswitch.Value
			updateHost = append(updateHost, v.host)
		}

		if v.saveBtn != nil && v.saveBtn.Clicked() {
			editorSaved = true
			v.host.Content = v.editor.Text()
			updateHost = append(updateHost, v.host)
		}

	}

	var needReload bool
	var itemdelete bool
	// 删除
	for m.mins.Clicked() {
		if len(m.tabs.tabs) <= 1 {
			break
		}
		if len(m.tabs.tabs) < m.tabs.selected {
			break
		}
		err := m.hostsManager.Delete(m.tabs.tabs[m.tabs.selected].host)
		if err != nil {
			break
		}
		itemdelete = true
		needReload = true
	}

	// 添加
	for m.plus.Clicked() {
		title, _, err := dlgs.Entry("Add hosts", "please name a hosts title", "")
		if err != nil {
			break
		}
		if title == "" {
			break
		}
		err = m.hostsManager.Create(&model.Host{Title: title})
		if err != nil {
			break
		}
		needReload = true
	}

	for _, v := range updateHost {
		m.hostsManager.Update(v)
	}

	// 有两个事件触发 写hosts文件 1. switch 改变 2. switch开的editor ctrl + s
	// 添加和删除
	if switchChanged || editorSaved || itemdelete {
		//m.hfm.Write(changeText
		if m.hostsManager.GetPwd() == "" && runtime.GOOS != "windows" {

			pwd, t, _ := dlgs.Password("password", "")
			if t {
				m.hostsManager.SetPwd(pwd)
			}
		}
		m.hostsManager.UpdateHostFile()
	}

	if needReload {
		m.invalidate <- struct{}{}
	}

}
