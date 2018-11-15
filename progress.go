package progress

import (
	ui "github.com/gizak/termui"
	"strings"
)

type Progress struct {
	border bool // objects have borders

	gs []*ui.Gauge // list of progress bars
	ps []*ui.Par   // list of associated info

	header *ui.Par // a "header" text box

	mtext    *ui.Par  // a "messages" text box
	messages []string // list of messages to return
}

type newOption func(p *Progress)

// Messages specify that we want a "messages" section with "m" lines
func Messages(m int) newOption {
	return func(p *Progress) {
		if p.border {
			m += 2
		}
		p.mtext = ui.NewPar("")
		p.mtext.Border = p.border
		p.mtext.Height = m
	}
}

// Header specify that we want a "header" section with "h" lines
func Header(h int) newOption {
	return func(p *Progress) {
		if p.border {
			h += 2
		}
		p.header = ui.NewPar("")
		p.header.Border = p.border
		p.header.Height = h
	}
}

//
// New creates a new Progress object that manages "n" progress item and a message area with "m" lines
//
func New(n int, border bool, options ...newOption) *Progress {
	p := &Progress{border: border, gs: make([]*ui.Gauge, n)}

	height, padding := 2, 1

	for _, opt := range options {
		opt(p) // process option
	}

	if p.border {
		height, padding = 3, 0
	} else {
		p.ps = make([]*ui.Par, n)
	}

	if p.header != nil {
		ui.Body.AddRows(ui.NewRow(ui.NewCol(12, 0, p.header)))
	}

	for i := 0; i < n; i++ {
		p.gs[i] = ui.NewGauge()
		p.gs[i].LabelAlign = ui.AlignLeft
		p.gs[i].Height = height
		p.gs[i].Border = p.border
		p.gs[i].Percent = 0
		p.gs[i].PaddingBottom = padding
		p.gs[i].BarColor = ui.ColorGreen

		if p.border {
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

	if p.mtext != nil {
		ui.Body.AddRows(ui.NewRow(ui.NewCol(12, 0, p.mtext)))
	}

	ui.Body.Align()
	ui.Render(ui.Body)

	return p
}

func PercInt(curr, max int) int {
	if curr == 0 || max == 0 {
		return 0
	}

	return curr * 100 / max
}

func PercInt64(curr, max int64) int {
	if curr == 0 || max == 0 {
		return 0
	}

	return int(curr * 100 / max)
}

func PercFloat(curr, max float64) int {
	if curr == 0.0 || max == 0.0 {
		return 0
	}

	return int(curr * 100.0 / max)
}

func (p *Progress) SetHeader(m string) {
	if p.header != nil {
		p.header.Text = m
	}
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

func (p *Progress) AddMessage(m string) {
	p.messages = append(p.messages, m)

	if p.mtext != nil {
		start := len(p.messages) - p.mtext.Height
		if start < 0 {
			start = 0
		}

		p.mtext.Text = strings.Join(p.messages[start:], "\n")
		ui.Render(ui.Body)
	}
}

func (p *Progress) Messages() []string {
	return p.messages
}
