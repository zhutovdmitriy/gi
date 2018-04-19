// Copyright (c) 2018, Randall C. O'Reilly. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gi

import (
	"fmt"

	"github.com/rcoreilly/goki/gi/units"
	"github.com/rcoreilly/goki/ki"
)

// open an interactive editor of the given Ki tree, at its root
func GoGiEditorOf(obj ki.Ki) *Window {
	width := 1280
	height := 1024
	win := NewWindow2D("GoGi Editor Window", width, height)
	win.UpdateStart()

	vp := win.WinViewport2D()
	vp.SetProp("background-color", "#FFF")
	vp.Fill = true

	vlay := vp.AddNewChild(KiT_Frame, "vlay").(*Frame)
	vlay.Lay = LayoutCol

	trow := vlay.AddNewChild(KiT_Layout, "trow").(*Layout)
	trow.Lay = LayoutRow
	trow.SetProp("align-vert", AlignMiddle)
	trow.SetProp("align-horiz", "center")
	trow.SetProp("margin", 2.0) // raw numbers = px = 96 dpi pixels
	trow.SetStretchMaxWidth()

	spc := vlay.AddNewChild(KiT_Space, "spc1").(*Space)
	spc.SetFixedHeight(units.NewValue(2.0, units.Em))

	trow.AddNewChild(KiT_Stretch, "str1")
	title := trow.AddNewChild(KiT_Label, "title").(*Label)
	title.Text = fmt.Sprintf("GoGi Editor of Ki Node Tree: %v", obj.Name())
	title.SetProp("max-width", -1)
	title.SetProp("text-align", "center")
	trow.AddNewChild(KiT_Stretch, "str2")

	split := vlay.AddNewChild(KiT_SplitView, "split").(*SplitView)
	split.Dim = X

	tvfr := split.AddNewChild(KiT_Frame, "tvfr").(*Frame)
	svfr := split.AddNewChild(KiT_Frame, "svfr").(*Frame)
	split.SetSplits(.3, .7)

	tv := tvfr.AddNewChild(KiT_TreeView, "tv").(*TreeView)
	tv.SetSrcNode(obj)

	sv := svfr.AddNewChild(KiT_StructView, "sv").(*StructView)
	sv.SetStruct(obj, nil)

	tv.TreeViewSig.Connect(sv.This, func(recv, send ki.Ki, sig int64, data interface{}) {
		if data == nil {
			return
		}
		tvn, _ := data.(ki.Ki).EmbeddedStruct(KiT_TreeView).(*TreeView)
		svr, _ := recv.EmbeddedStruct(KiT_StructView).(*StructView)
		if sig == int64(TreeViewSelected) {
			svr.SetStruct(tvn.SrcNode.Ptr, nil)
		}
	})

	bspc := vlay.AddNewChild(KiT_Space, "ButSpc").(*Space)
	bspc.SetFixedHeight(units.NewValue(1.0, units.Em))

	brow := vlay.AddNewChild(KiT_Layout, "brow").(*Layout)
	brow.Lay = LayoutRow
	brow.SetProp("align-vert", AlignMiddle)
	brow.SetProp("align-horiz", "center")
	brow.SetProp("margin", 2.0) // raw numbers = px = 96 dpi pixels
	brow.SetStretchMaxWidth()

	updtobj := brow.AddNewChild(KiT_Button, "updtobj").(*Button)
	updtobj.SetText("Update")
	updtobj.ButtonSig.Connect(win.This, func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(ButtonClicked) {
			obj.UpdateStart()
			obj.UpdateEnd()
		}
	})

	savej := brow.AddNewChild(KiT_Button, "savejson").(*Button)
	savej.SetText("Save JSON")
	savej.ButtonSig.Connect(win.This, func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(ButtonClicked) {
			obj.SaveJSONToFile("GoGiEditorOut.json") // todo: first a string prompt, then a file dialog
		}
	})

	loadj := brow.AddNewChild(KiT_Button, "loadjson").(*Button)
	loadj.SetText("Load JSON")
	loadj.ButtonSig.Connect(win.This, func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(ButtonClicked) {
			obj.LoadJSONFromFile("GoGiEditorOut.json") // todo: first a string prompt, then a file dialog
		}
	})

	win.UpdateEnd()
	win.StartEventLoopNoWait()
	return win
}