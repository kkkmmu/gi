// Copyright (c) 2018, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package giv

import (
	"fmt"

	"github.com/goki/gi"
	"github.com/goki/gi/oswin"
	"github.com/goki/ki"
)

// PrefsView opens a view of user preferences
func PrefsView(pf *gi.Preferences) (*StructView, *gi.Window) {
	winm := "gogi-prefs"
	if w, ok := gi.MainWindows.FindName(winm); ok {
		w.OSWin.Raise()
		return nil, nil
	}

	width := 800
	height := 800
	win := gi.NewWindow2D(winm, "GoGi Preferences", width, height, true)

	vp := win.WinViewport2D()
	updt := vp.UpdateStart()

	mfr := win.SetMainFrame()
	mfr.Lay = gi.LayoutVert

	sv := mfr.AddNewChild(KiT_StructView, "sv").(*StructView)
	sv.Viewport = vp
	sv.SetStruct(pf, nil)
	sv.SetStretchMaxWidth()
	sv.SetStretchMaxHeight()

	mmen := win.MainMenu
	MainMenuView(pf, win, mmen)

	inClosePrompt := false
	win.OSWin.SetCloseReqFunc(func(w oswin.Window) {
		if pf.Changed {
			if !inClosePrompt {
				gi.ChoiceDialog(vp, gi.DlgOpts{Title: "Save Prefs Before Closing?",
					Prompt: "Do you want to save any changes to preferences before closing?"},
					[]string{"Save and Close", "Discard and Close", "Cancel"},
					win.This, func(recv, send ki.Ki, sig int64, data interface{}) {
						switch sig {
						case 0:
							pf.Save()
							fmt.Println("Preferences Saved to prefs.json")
							w.Close()
						case 1:
							pf.Open() // if we don't do this, then it actually remains in edited state
							w.Close()
						case 2:
							inClosePrompt = false
							// default is to do nothing, i.e., cancel
						}
					})
			}
		} else {
			w.Close()
		}
	})

	win.MainMenuUpdated()

	vp.UpdateEndNoSig(updt)
	win.GoStartEventLoop()
	return sv, win
}
