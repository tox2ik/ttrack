package tui

import (
	"github.com/jroimartin/gocui"
)


var wrapEnabled = false

func cursorDown(g *gocui.Gui, v *gocui.View) (err error) {
	err = nil
	if nil == v {
		return
	}

	cx, cy := v.Cursor()
	if cy < ui.Len()-1 {
		err = v.SetCursor(cx, cy+1)
	} else {
		if wrapEnabled {
			err = v.SetCursor(cx, 0)
		}
	}

	if err != nil {
		ox, oy := v.Origin()
		if err = v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}

	return
}

func cursorUp(gui *gocui.Gui, view *gocui.View) error {
	if view != nil {
		ox, oy := view.Origin()
		cx, cy := view.Cursor()
		// defer func(){
		// 	_, _oy := view.Origin()
		// 	_, _cy := view.Cursor()
		// 	debug("origin %d cursor %d {up}\n", _oy, _cy)
		// }()

		err := view.SetCursor(cx, cy-1)

		if err != nil && oy > 0 {
			if err := view.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}
