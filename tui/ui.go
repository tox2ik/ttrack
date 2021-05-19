package tui

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/nsf/termbox-go"

	"genja.org/ttrack/model"
	"genja.org/ttrack/ttio"
)

var ui *UiState

const (
	ViewNone   = ``
	ViewStamps = `stamps`
	ViewDebug  = `debug`
	ViewAsk    = `ask`
	ViewAskInp = `ask_input`
)

func Run(arg model.Arguments) (err error) {
	var gui *gocui.Gui

	gui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer gui.Close()

	gui.SetManagerFunc(layout)
	err = bindKeys(gui)
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(nil)
	initState(arg, gui, buf)

	go func() {
		time.Sleep(time.Millisecond * 133)
		gui.Update(func(gui *gocui.Gui) error {
			gw := guiEw(gui)
			gw.Write(buf.Bytes())
			return nil
		})
	}()

	err = gui.MainLoop()
	if err == gocui.ErrQuit {
		return nil
	}
	return
}

func initState(arg model.Arguments, gui *gocui.Gui, ew io.Writer) {
	if nil == ui {
		ui = &UiState{
			Args:         arg,
			Gui:          gui,
			Presentation: PresiStamps{},
			DebugVisible: true,
		}
	}
	rec, tup, _ := ttio.ParseRecords(arg, ew)
	ui.Stamps = tup
	ui.Records = rec
}

func redraw(g *gocui.Gui) {
	vStamps, _ := g.View(ViewStamps)
	if vStamps == nil {
		return
	}
	ocx, ocy := vStamps.Cursor()
	vStamps.Clear()
	ui.Draw(vStamps)
	if ocy == ui.Len() {
		vStamps.SetCursor(ocx, ocy-1)
	} else
	if ocy > ui.Len() {
		vStamps.SetCursor(ocx, ui.Len()-1)
	} else {
		vStamps.SetCursor(ocx, ocy)
	}
}

func layout(g *gocui.Gui) (err error) {
	var v *gocui.View

	debugH := 0
	if ui.DebugVisible {
		debugH = 7
	}

	X, Y := g.Size()
	v, err = g.SetView(ViewStamps, 0, 0, X-1, Y-debugH-1)

	if err == gocui.ErrUnknownView && v != nil {
		ui.StampView = v
		v.Highlight = true
		v.SelFgColor = gocui.Attribute(termbox.ColorLightGray)
		v.SelBgColor = gocui.ColorBlack
		_, err = g.SetCurrentView(ViewStamps)
		cy := ui.Len() - 1
		_ = v.SetCursor(0, ui.inboundsY(cy))
		redraw(g)
	}

	if ui.DebugVisible {
		v, err = g.SetView(ViewDebug, 0, Y-debugH-1, X-1, Y-1)
		if err == gocui.ErrUnknownView && v != nil {
			v.Frame = false
			v.Autoscroll = true
			err = nil
		}
	}

	return
}

func debug(f string, a ...interface{}) {
	// v, _ := ui.Gui.View("debug")
	// //_, y := v.Size()
	// //if len(v.BufferLines()) >= y { v.Clear() }
	// if v != nil {
	// }
	fmt.Fprintf(guiEw(ui.Gui), f, a...)
}

func guiEw(gui *gocui.Gui) io.Writer {
	var wr io.Writer
	var err error
	wr, err = gui.View(ViewDebug)
	if err != nil {
		wr = os.Stderr
	}
	return wr
}

// quit is invoked when the user presses "Ctrl+C"
func quit(*gocui.Gui, *gocui.View) error {
	return gocui.ErrQuit
}

// globalQuit is invoked when the user quits the contact and or
// when all conflicts have been resolved
func globalQuit(g *gocui.Gui, err error) {
	g.Update(func(g *gocui.Gui) error {
		return err
	})
}
