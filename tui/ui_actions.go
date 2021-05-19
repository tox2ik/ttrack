package tui

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/nsf/termbox-go"

	"genja.org/ttrack/ttio"
)


var stampPresentations = []Presi{
	PresiStamps{},
	PresiRecords{},
	PresiPerDay{},
}
var currentPresentation = 0

func bindKeys(g *gocui.Gui) (err error) {
	resp := append([]error{},
		g.SetKeybinding("", gocui.KeyCtrlL, gocui.ModNone, ctrlL),
		nil,
		g.SetKeybinding(ViewNone, 'q', gocui.ModNone, quit),
		g.SetKeybinding(ViewNone, 'v', gocui.ModNone, togglePresentation),
		g.SetKeybinding(ViewNone, 'c', gocui.ModNone, toggleDebug),
		g.SetKeybinding(ViewNone, gocui.KeyCtrlC, gocui.ModNone, quit),
		nil,
		g.SetKeybinding(ViewAskInp, gocui.KeyEnter, gocui.ModNone, parseInput),
		nil,
		g.SetKeybinding(ViewStamps, gocui.KeyArrowDown, gocui.ModNone, cursorDown),
		g.SetKeybinding(ViewStamps, gocui.KeyArrowUp, gocui.ModNone, cursorUp),
		g.SetKeybinding(ViewStamps, gocui.KeyEnter, gocui.ModNone, addStamp),
		g.SetKeybinding(ViewStamps, gocui.KeyEnter, gocui.ModAlt, selectItem),
		g.SetKeybinding(ViewStamps, 'a', gocui.ModNone, askStamp),
		g.SetKeybinding(ViewStamps, 'j', gocui.ModNone, cursorDown),
		g.SetKeybinding(ViewStamps, 'k', gocui.ModNone, cursorUp),
		g.SetKeybinding(ViewStamps, 'd', gocui.ModNone, markDelete),
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


func toggleDebug(gui *gocui.Gui, view *gocui.View) error {
	ui.DebugVisible = ! ui.DebugVisible

	if ui.DebugVisible {
		gui.DeleteView(ViewDebug)
	}
	redraw(gui)
	return nil
}

func togglePresentation(gui *gocui.Gui, view *gocui.View) error {
	currentPresentation++
	if currentPresentation >= len(stampPresentations) {
		currentPresentation = 0
	}
	ui.Presentation = stampPresentations[currentPresentation % len(stampPresentations) ]
	redraw(gui)
	return nil
}

func ctrlL(*gocui.Gui, *gocui.View) error {
	return termbox.Sync()
}


func askStamp(gui *gocui.Gui, view *gocui.View) (err error) {

	x, y := gui.Size()
	H := 1
	W := 20

	v, err := gui.SetView(ViewAsk, x/3-W, y/4-1-H, x/3+W, y/4-1+H)
	if err == gocui.ErrUnknownView {
		fmt.Fprint(v, "Enter time:")
	}

	vi, err := gui.SetView(ViewAskInp, x/3-W+12, y/4-1-H, x/3+W, y/4-1+H)
	if err == gocui.ErrUnknownView {
		vi.Editable = true
		vi.Frame = false
		_, _ = gui.SetCurrentView(ViewAskInp);
		gui.SetViewOnTop(ViewAskInp)

	} else  {
		debug("ask: %s\n", err)
	}

	return nil
}

func parseInput(gui *gocui.Gui, view *gocui.View) error {

	v, _ := gui.View(ViewAskInp)
	if v != nil {
		t, err := ttio.ParseDate(v.Buffer())
		if err != nil {

			debug("%s", err)
		} else {
			v.Clear()
			gui.SetCurrentView(ViewStamps)
			gui.DeleteView(ViewAskInp)
			gui.DeleteView(ViewAsk)

			ui.Args.Stamp = t
			ttio.AddStamp(ui.Args, guiEw(gui))
			initState(ui.Args, gui, guiEw(gui))
			redraw(gui)
		}




	}
	return nil
}

func addStamp(gui *gocui.Gui, view *gocui.View) (err error) {
	ttio.AddStamp(ui.Args, guiEw(gui))
	initState(ui.Args, gui, guiEw(gui))
	redraw(gui)
	return
}


func markDelete(gui *gocui.Gui, view *gocui.View) error {
	if view != nil {
		ui.RemoveItem(view)

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

