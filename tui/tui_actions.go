package tui

import (
	"io"
	"os"

	"github.com/jroimartin/gocui"
	"github.com/nsf/termbox-go"

	"genja.org/ttrack/ttio"
)

func bindKeys(g *gocui.Gui) (err error) {
	resp := append([]error{},
		g.SetKeybinding("", gocui.KeyCtrlL, gocui.ModNone, ctrlL),
		nil,
		g.SetKeybinding(ViewNone, 'q', gocui.ModNone, quit),
		g.SetKeybinding(ViewNone, gocui.KeyCtrlC, gocui.ModNone, quit),
		nil,
		g.SetKeybinding(ViewStamps, gocui.KeyArrowDown, gocui.ModNone, cursorDown),
		g.SetKeybinding(ViewStamps, gocui.KeyArrowUp, gocui.ModNone, cursorUp),
		g.SetKeybinding(ViewStamps, gocui.KeyEnter, gocui.ModNone, addStamp),
		g.SetKeybinding(ViewStamps, gocui.KeyEnter, gocui.ModAlt, selectItem),
		g.SetKeybinding(ViewStamps, 'j', gocui.ModNone, cursorDown),
		g.SetKeybinding(ViewStamps, 'k', gocui.ModNone, cursorUp),
		g.SetKeybinding(ViewStamps, 'd', gocui.ModNone, deleteItem),
		g.SetKeybinding(ViewStamps, 'w', gocui.ModNone, flushItems),
		g.SetKeybinding(ViewStamps, 's', gocui.ModNone, swapRecords),
	)
	for _, e := range resp {
		if e != nil {
			err = e
		}
	}
	return
}

func ctrlL(*gocui.Gui, *gocui.View) error {
	return termbox.Sync()
}

func guiEw(gui *gocui.Gui) io.Writer {
	var wr io.Writer
	var err error
	wr, err = gui.View(ViewDebug)
	if err != nil {
		wr = os.Stderr
	}
	return  wr
}

func addStamp(gui *gocui.Gui, view *gocui.View) (err error) {
	ttio.AddStamp(ui.Args, guiEw(gui))
	initState(ui.Args, gui, guiEw(gui))
	redraw(gui)
	//ctrlL(nil, nil)
	return
}

var WrapEnabled = false

func cursorDown(g *gocui.Gui, v *gocui.View) (err error) {
	err = nil
	if nil == v {
		return
	}

	cx, cy := v.Cursor()
	if cy < ui.Len()-1 {
		err = v.SetCursor(cx, cy+1)
	} else {
		if WrapEnabled {
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
		if err := view.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := view.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func deleteItem(gui *gocui.Gui, view *gocui.View) error {
	if view != nil {
		_, cy := view.Cursor()
		tup := ui.RemoveStamp(cy)
		ui.RemovedStamps = append(ui.RemovedStamps, tup)
		redraw(gui)
	}
	return nil
}

func flushItems(gui *gocui.Gui, view *gocui.View) error {
	for _, rms := range ui.RemovedStamps {
		ttio.RemoveFileRecord(ui.Args, rms.In)
		ttio.RemoveFileRecord(ui.Args, rms.Out)
	}
	for _, rm := range ui.RemovedRecords {
		ttio.RemoveFileRecord(ui.Args, rm)
	}
	return nil
}

func swapRecords(gui *gocui.Gui, view *gocui.View) error {
	if view != nil {
		_, cy := view.Cursor()
		ui.SwapRecords(cy)
		redraw(gui)
	}
	return nil

}

func selectItem(gui *gocui.Gui, view *gocui.View) error {
	if view == nil {
		return nil
	}

	_, cy := view.Cursor()
	i := ui.inboundsY(cy)
	tup := ui.Stamps.Items[i]
	debug("\r" + tup.FormatHuman())
	redraw(gui)

	return nil

}

const (
	ViewNone    = ``
	ViewStamps  = `stamps`
	ViewRecords = `records`
	ViewDebug   = `debug`
)
