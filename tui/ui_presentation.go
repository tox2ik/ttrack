package tui

import (
	"github.com/jroimartin/gocui"

	"genja.org/ttrack/glue"
	"genja.org/ttrack/model"
)

type Presi interface {
	Draw(v *gocui.View)
	Len(stamps model.Tuples) int
}

type (
	PresiStamps struct{}
	PresiPerDay struct{}
	PresiRecords struct{}
)

func (p PresiStamps) Draw(v *gocui.View) {
	err := ui.Stamps.ReportHoursHuman(v)
	glue.Die(err)
}

func (p PresiPerDay) Draw(v *gocui.View) {
	err := ui.Stamps.ReportHoursPerDay(v)
	glue.Die(err)
}

func (p PresiRecords) Draw(v *gocui.View) {
	err := ui.Stamps.ReportRecords(v)
	glue.Die(err)
}

func (p PresiStamps) Len(stamps model.Tuples) int {
	return len(stamps.Items)
}

func (p PresiPerDay) Len(stamps model.Tuples) int {
	d, _ := stamps.SecondsPerDay()
	return len(d)
}
func (p PresiRecords) Len(stamps model.Tuples) int {
	c := 0
	for _, t := range stamps.Items {
		if t.In.IsValid() {
			c++
		}
		if t.Out.IsValid() {
			c++
		}
	}
	return c
}
