package tui

import (
	"github.com/jroimartin/gocui"

	"genja.org/ttrack/model"
)

type UiState struct {
	Args           model.Arguments
	Stamps         model.Tuples
	Records        []model.Record
	RemovedStamps  []model.Tuple
	RemovedRecords []model.Record
	Presentation   Presi
	DebugVisible   bool

	Gui            *gocui.Gui
	StampView      *gocui.View
}

func (ui UiState) Len() int {

	//return len(ui.Stamps.Items)
	return ui.Presentation.Len(ui.Stamps)
}

func (ui UiState) inboundsY(y int) int {
	i := y
	if i < 0 {
		i = 0
	}
	if i >= ui.Len() {
		i = ui.Len() - 1
	}
	return i

}

func (ui *UiState) RemoveItem(view *gocui.View) (*model.Tuple, *model.Record) {




	if _, kk := ui.Presentation.(PresiStamps); kk {
		_, oy := view.Origin()
		_, cy := view.Cursor()
		y := oy + cy
		removed := ui.Stamps.Remove(y)
		ui.RemovedStamps = append(ui.RemovedStamps, removed)
		return &removed, nil
	}

	if _, kk := ui.Presentation.(PresiRecords); kk {

		_, oy := view.Origin()
		_, cy := view.Cursor()
		y := oy + cy

		tupi := y / 2
		reci := y % 2

		tup := &ui.Stamps.Items[tupi]
		rmrec := model.Record{}

		if reci == 0 {
			rmrec = tup.In
			tup.In = model.Record{
				Mark:  "in",
				Day:   tup.In.Day,
				Time:  "---",
				Stamp: 0,
			}
		}
		if reci == 1 {
			rmrec = tup.Out
			tup.Out = model.Record{
				Mark:  "out",
				Day:   tup.In.Day,
				Time:  ">>>",
				Stamp: 0,
			}
		}
		debug("rec: %s\n", rmrec.String())

		ui.RemovedRecords = append(ui.RemovedRecords, rmrec)
		if tup.Len() == 0 {
			rmt := ui.Stamps.Remove(cy)
			ui.RemovedStamps = append(ui.RemovedStamps, rmt)

		}
		return nil, nil
	}

	return nil, nil

}

func (ui *UiState) SwapRecords(cy int) {
	ui.Stamps.Items[ui.inboundsY(cy)].Swap()
}

func (ui UiState) Draw(v *gocui.View) {
	ui.Presentation.Draw(v)
}
