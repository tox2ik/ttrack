package tui

import (
	"fmt"
	"io"

	"github.com/jroimartin/gocui"
	"github.com/nsf/termbox-go"

	"genja.org/ttrack/model"
	"genja.org/ttrack/ttio"
)

var ui *UiState

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

	initState(arg, gui, guiEw(gui))

	err = gui.MainLoop()
	if err == gocui.ErrQuit {
		return nil
	}
	return
}

func initState(args model.Arguments, gui *gocui.Gui, ew io.Writer) {
	if nil == ui {
		ui = &UiState{Args: args, Gui: gui}
	}
	rec, tup, _ := ttio.ParseRecords(args, ew)
	ui.Stamps = tup
	ui.Records = rec
}

func redraw(g *gocui.Gui) {

	vStamps, _ := g.View(ViewStamps)
	if vStamps != nil {
		ocx, ocy := vStamps.Cursor()

		vStamps.Clear()
		for _, t := range ui.Stamps.Items {
			fmt.Fprintf(vStamps, t.FormatHuman())
		}

		fmt.Fprintf(vStamps, "                      total: %s\n", ui.Stamps.HoursH())
		fmt.Fprintf(vStamps, "                    average: %s\n", ui.Stamps.HoursAverageH())

		if ocy == ui.Len() {
			vStamps.SetCursor(ocx, ocy-1)
		} else {
			vStamps.SetCursor(ocx, ocy)
		}
	}
}

func layout(g *gocui.Gui) (err error) {
	var v *gocui.View

	X, Y := g.Size()
	v, err = g.SetView(ViewStamps, 0, 0, X-1, Y-1)

	if err == gocui.ErrUnknownView && v != nil {
		v.Highlight = true
		v.SelFgColor = gocui.Attribute(termbox.ColorLightGray)
		v.SelBgColor = gocui.ColorBlack
		_, err = g.SetCurrentView(ViewStamps)
		cy := ui.Len() - 1
		_ = v.SetCursor(0, ui.inboundsY(cy))
		redraw(g)
	}


	v, err = g.SetView(ViewDebug, 1, Y-7, X-2, Y-2)
	if err == gocui.ErrUnknownView && v != nil {
		err = nil
	}
	return
}

func debug(f string, a ...interface{}) {
	v, _ := ui.Gui.View("debug")
	_, y := v.Size()
	if len(v.BufferLines()) >= y {
		v.Clear()
	}
	if v != nil {
		fmt.Fprintf(v, f, a...)
	}
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
