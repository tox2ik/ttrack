package tui

import (
	"github.com/jroimartin/gocui"

	"genja.org/ttrack/model"
)

type UiState struct {
	Args    model.Arguments
	Stamps  model.Tuples
	Records []model.Record
	Gui     *gocui.Gui
	RemovedStamps []model.Tuple
	RemovedRecords []model.Record
}

func (ui UiState) Len() int {
	return len(ui.Stamps.Items)
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

func (ui *UiState) RemoveStamp(cy int) model.Tuple {
	removed := ui.Stamps.Remove(cy)
	ui.RemovedStamps = append(ui.RemovedStamps, removed)
	return removed
}

func (ui *UiState) SwapRecords(cy int) {
	ui.Stamps.Items[ui.inboundsY(cy)].Swap()
}
