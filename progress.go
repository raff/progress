package progress

import (
	ui "github.com/gizak/termui"
)

type Progress struct {
	border bool
	gs     []*ui.Gauge
	ps     []*ui.Par
}

func New(n int, border bool) *Progress {
	p := Progress{border: border}
	p.gs = make([]*ui.Gauge, n)

	height, padding := 2, 1

	if border {
		height, padding = 3, 0
	} else {
		p.ps = make([]*ui.Par, n)
	}

	for i := 0; i < n; i++ {
		p.gs[i] = ui.NewGauge()
		p.gs[i].LabelAlign = ui.AlignLeft
		p.gs[i].Height = height
		p.gs[i].Border = border
		p.gs[i].Percent = 0
		p.gs[i].PaddingBottom = padding
		p.gs[i].BarColor = ui.ColorGreen

		if border {
			ui.Body.AddRows(ui.NewRow(ui.NewCol(12, 0, p.gs[i])))
		} else {
			p.ps[i] = ui.NewPar("")
			p.ps[i].Height = 1
			p.ps[i].Border = false

			// build layout
			ui.Body.AddRows(ui.NewRow(
				ui.NewCol(4, 0, p.ps[i]),
				ui.NewCol(8, 0, p.gs[i])))
		}
	}

	ui.Body.Align()
	ui.Render(ui.Body)

	return &p
}

func PercInt(curr, max int) int {
	if curr == 0 || max == 0 {
		return 0
	}

	return curr * 100 / max
}

func PercFloat(curr, max float64) int {
	if curr == 0.0 || max == 0.0 {
		return 0
	}

	return int(curr * 100.0 / max)
}

func (p *Progress) Set(item int, label string, value int) {
	if p.border {
		p.gs[item].BorderLabel = label
	} else {
		p.ps[item].Text = label
	}
	p.gs[item].Percent = value
	ui.Render(ui.Body)
}
